package aws

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Manager struct {
	meta struct {
		key    string
		addr   string
		secret string

		concurrency int
	}

	cli *s3.Client

	uploader   *manager.Uploader
	downloader *manager.Downloader
}

func newDefaultS3Manager() *S3Manager {
	s3m := &S3Manager{}
	s3m.meta.concurrency = 3
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

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3m.meta.key,
				s3m.meta.secret,
				"",
			)),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: s3m.meta.addr,
					}, nil
				})),
	)
	if err != nil {
		return nil, err
	}

	s3m.cli = s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	s3m.uploader = manager.NewUploader(
		s3m.cli,
		func(uploader *manager.Uploader) {
			uploader.Concurrency = s3m.meta.concurrency
		},
	)
	s3m.downloader = manager.NewDownloader(
		s3m.cli,
		func(downloader *manager.Downloader) {
			downloader.Concurrency = s3m.meta.concurrency
		},
	)

	return s3m, nil
}

func (s *S3Manager) createBucket(ctx context.Context, bucket string) error {
	_, err := s.cli.CreateBucket(
		ctx,
		&s3.CreateBucketInput{
			Bucket: &bucket,
		},
	)
	return err
}

func (s *S3Manager) Upload(ctx context.Context, bucket, objectKey string, data []byte) error {
	_, err := s.uploader.Upload(
		ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
			Body:   bytes.NewReader(data),
		},
	)
	return err
}

func (s *S3Manager) Download(ctx context.Context, bucket, objectKey string) ([]byte, error) {
	headObject, err := s.cli.HeadObject(
		ctx,
		&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		},
	)

	buf := make([]byte, int(headObject.ContentLength))
	_, err = s.downloader.Download(
		ctx,
		manager.NewWriteAtBuffer(buf),
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(objectKey),
		})

	return buf, err
}
