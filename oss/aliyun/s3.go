package aliyun

import (
	"bytes"
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
)

type S3Manager struct {
	meta struct {
		addr   string
		key    string
		secret string
	}

	cli *oss.Client
}

func newDefaultS3Manager() *S3Manager {
	s3m := &S3Manager{}
	return s3m
}

func NewS3Cli(opts ...Option) (*S3Manager, error) {
	s3m := newDefaultS3Manager()

	for _, opt := range opts {
		opt.apply(s3m)
	}
	if len(s3m.meta.key) == 0 {
		return nil, errWrapper(IllegalParams)
	}
	if len(s3m.meta.secret) == 0 {
		return nil, errWrapper(IllegalParams)
	}
	if len(s3m.meta.addr) == 0 {
		return nil, errWrapper(IllegalParams)
	}

	cli, err := oss.New(s3m.meta.addr, s3m.meta.key, s3m.meta.secret)
	if err != nil {
		return nil, err
	}
	s3m.cli = cli

	return s3m, nil
}

func (s *S3Manager) CreateBucket(_ context.Context, bucket string) error {
	err := s.cli.CreateBucket(
		bucket,
	)
	return err
}

func (s *S3Manager) Upload(_ context.Context, bucket, objectKey string, data []byte) error {
	bk, err := s.cli.Bucket(bucket)
	if err != nil {
		return err
	}
	return bk.PutObject(objectKey, bytes.NewReader(data))
}

func (s *S3Manager) Download(_ context.Context, bucket, objectKey string) ([]byte, error) {
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
