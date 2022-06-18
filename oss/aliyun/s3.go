package aliyun

import (
	"bytes"
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
)

type OssManager struct {
	meta struct {
		addr   string
		key    string
		secret string
	}

	cli *oss.Client
}

func newDefaultOssManager() *OssManager {
	return &OssManager{}
}

func NewOssCli(opts ...Option) (*OssManager, error) {
	ossM := newDefaultOssManager()

	{
		for _, opt := range opts {
			opt.apply(ossM)
		}
		switch {
		case len(ossM.meta.key) == 0:
			return nil, errWrapper(IllegalParams)
		case len(ossM.meta.secret) == 0:
			return nil, errWrapper(IllegalParams)
		case len(ossM.meta.addr) == 0:
			return nil, errWrapper(IllegalParams)
		}
	}

	cli, err := oss.New(ossM.meta.addr, ossM.meta.key, ossM.meta.secret)
	if err != nil {
		return nil, err
	}
	ossM.cli = cli

	return ossM, nil
}

func (s *OssManager) CreateBucket(_ context.Context, bucket string) error {
	err := s.cli.CreateBucket(
		bucket,
	)
	return err
}

func (s *OssManager) Upload(_ context.Context, bucket, objectKey string, data []byte) error {
	bk, err := s.cli.Bucket(bucket)
	if err != nil {
		return err
	}
	return bk.PutObject(objectKey, bytes.NewReader(data))
}

func (s *OssManager) Download(_ context.Context, bucket, objectKey string) ([]byte, error) {
	bk, err := s.cli.Bucket(bucket)
	if err != nil {
		return nil, err
	}
	resp, err := bk.GetObject(objectKey)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	return ioutil.ReadAll(resp)
}
