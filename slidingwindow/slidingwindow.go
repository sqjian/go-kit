package slidingwindow

import (
	"sync/atomic"
	"time"
)

type Bucket interface {
	Update(...any)
	Clear()
}

type SlidingWindow struct {
	//循环队列
	buckets []Bucket

	//队列的总长度
	timeSliceSize int64

	//每个时间片的时长
	timePerSlice time.Duration

	//窗口长度
	winSize int64

	//当前所使用的时间片位置
	cursor int64

	//上一次访问的时间戳
	preTs time.Time

	//窗口总时长
	winDur time.Duration
}

func NewSlidingWindow(timePerSlice time.Duration, winSize int64, bucketProvider func() Bucket) *SlidingWindow {
	inst := SlidingWindow{
		timePerSlice:  timePerSlice,
		winSize:       winSize,
		timeSliceSize: winSize*2 + 1,

		buckets: func() []Bucket {
			var tmp []Bucket
			for i := int64(0); i < winSize*2+1; i++ {
				tmp = append(tmp, bucketProvider())
			}
			return tmp
		}(),
		preTs:  time.Now(),
		winDur: timePerSlice * time.Duration(winSize),
	}

	go inst.reset()

	return &inst
}

func (s *SlidingWindow) Clear() {
	for index := 0; index < len(s.buckets); index++ {
		s.buckets[index].Clear()
	}
}

// 清理闲时的节点数据
func (s *SlidingWindow) reset() {

	ticker := time.NewTicker(s.winDur)

	for {
		select {
		case <-ticker.C:
			{
				if time.Now().Sub(s.preTs) < s.winDur {
					//访问间隔没有超过所有区间
					continue
				}
				s.Clear()
			}
		}
	}
}

func (s *SlidingWindow) locationIndex() int64 {
	return (time.Now().UnixNano() / int64(s.timePerSlice)) % s.timeSliceSize
}

func (s *SlidingWindow) Update(values ...any) {
	var index = s.locationIndex()

	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, s.locationIndex())

	if oldCursor == index {
		// 在当前时间片里继续
		s.buckets[index].Update(values)
	} else {
		s.buckets[index].Update(values)
		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	s.preTs = time.Now()
}

func (s *SlidingWindow) GetBuckets() []Bucket {
	var index = s.locationIndex()

	// cursor不等于index，将cursor设置为index
	oldCursor := atomic.LoadInt64(&s.cursor)
	atomic.StoreInt64(&s.cursor, index)

	if oldCursor != index {
		// 可能有其他goroutine已经置过，问题不大
		s.buckets[index].Clear()

		// 清零，访问量不大时会有时间片跳跃的情况
		s.clearBetween(oldCursor, index)
	}

	var rst []Bucket
	for i := int64(0); i < s.winSize; i++ {
		rst = append(
			rst,
			s.buckets[(index-i+s.timeSliceSize)%s.timeSliceSize],
		)
	}
	return rst
}

// 将fromIndex~toIndex之间的时间片计数都清零
func (s *SlidingWindow) clearBetween(fromIndex, toIndex int64) {
	if time.Since(s.preTs) > s.winDur {
		// 解决极端情况下，当循环队列已经走了超过1个timeSliceSize以上，这里的清零并不能如期望的进行
		s.Clear()
		return
	}

	for index := (fromIndex + 1) % s.timeSliceSize; index != toIndex; index = (index + 1) % s.timeSliceSize {
		s.buckets[index].Clear()
	}
}
