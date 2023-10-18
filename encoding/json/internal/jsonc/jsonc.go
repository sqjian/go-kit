package jsonc

// 定义常量
const (
	ESCAPE   = 92 // 转义字符'\'
	QUOTE    = 34 // 双引号'"'
	NEWLINE  = 10 // 换行符'\n'
	ASTERISK = 42 // 星号'*'
	SLASH    = 47 // 斜线'/'
	HASH     = 35 // 井号'#'
	CARRIAGE = 13 // 回车符'\r'
)

// Translate 函数处理字节数组，主要是处理注释并返回新的字节数组
func Translate(s []byte) []byte {
	// 初始化变量
	var (
		i       int  // 新字节数组的索引
		quote   bool // 判断是否在双引号内部
		escaped bool // 判断当前字符是否被转义
	)
	j := make([]byte, len(s))
	comment := &commentData{} // 注释数据
	for _, ch := range s {
		// 处理转义字符
		if ch == ESCAPE || escaped {
			if !comment.started {
				j[i] = ch
				i++
			}
			escaped = !escaped
			continue
		}
		// 判断是否为双引号
		if ch == QUOTE {
			quote = !quote
		}
		// 处理单行注释结束
		if ch == NEWLINE || ch == CARRIAGE {
			if comment.isSingleLined {
				comment.stop()
			}

			if !comment.started {
				// 保留非注释里的换行，这里 \r、\n会被逐次追加
				j[i] = ch
				i++
			}
			continue
		}
		// 当在双引号内部且没有开始注释时，保留字符
		if quote && !comment.started {
			j[i] = ch
			i++
			continue
		}
		// 处理多行注释结束
		if comment.started {
			if ch == ASTERISK && !comment.isSingleLined {
				comment.canEnd = true
				continue
			}
			if comment.canEnd && ch == SLASH && !comment.isSingleLined {
				comment.stop()
				continue
			}
			comment.canEnd = false
			continue
		}
		// 判断是否可以开始注释
		if comment.canStart && (ch == ASTERISK || ch == SLASH) {
			comment.start(ch)
			continue
		}
		if ch == SLASH {
			comment.canStart = true
			continue
		}
		if ch == HASH {
			comment.start(ch)
			continue
		}
		j[i] = ch
		i++
	}
	return j[:i]
}

// commentData结构体用于存储注释相关的数据
type commentData struct {
	canStart      bool // 是否可以开始注释
	canEnd        bool // 是否可以结束注释
	started       bool // 注释是否已经开始
	isSingleLined bool // 是否是单行注释
	endLine       int  // 注释结束的行号
}

// stop方法用于结束注释
func (c *commentData) stop() {
	c.started = false
	c.canStart = false
}

// start方法用于开始注释
func (c *commentData) start(ch byte) {
	c.started = true
	c.isSingleLined = ch == SLASH || ch == HASH
}
