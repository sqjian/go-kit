package s3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	meta struct {
		key    string
		addr   string
		secret string
		debug  bool

		logger DefLogger

		concurrency int
	}

	cfg aws.Config
	cli *s3.Client

	uploader   *manager.Uploader
	downloader *manager.Downloader
}

func newDefaultS3Config() *S3 {
	cli := &S3{}
	cli.meta.concurrency = 3
	return cli
}

func newS3Cli(opts ...Option) (*S3, error) {
	cli := newDefaultS3Config()

	{
		for _, opt := range opts {
			opt.apply(cli)
		}
		if len(cli.meta.key) == 0 {
			return nil, ErrWrapper(IllegalParams)
		}
		if len(cli.meta.secret) == 0 {
			return nil, ErrWrapper(IllegalParams)
		}
		if len(cli.meta.addr) == 0 {
			return nil, ErrWrapper(IllegalParams)
		}
	}

	{
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					cli.meta.key,
					cli.meta.secret,
					"",
				)),
			config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "",
					SigningRegion: "",
					URL:           cli.meta.addr,
				}, nil
			})),
			config.WithLogger(&DefLogger{dummy: !cli.meta.debug}),
			config.WithClientLogMode(
				aws.LogSigning|
					aws.LogRetries|
					aws.LogRequest|
					aws.LogRequestWithBody|
					aws.LogResponse|
					aws.LogResponseWithBody|
					aws.LogDeprecatedUsage|
					aws.LogRequestEventMessage|
					aws.LogResponseEventMessage),
		)
		if err != nil {
			return nil, err
		}

		cli.cfg = cfg
	}

	{
		cli.cli = s3.NewFromConfig(cli.cfg, func(options *s3.Options) {
			options.UsePathStyle = true
		})
	}
	{
		cli.uploader = manager.NewUploader(
			cli.cli,
			func(uploader *manager.Uploader) {
				uploader.Concurrency = cli.meta.concurrency
			},
		)
		cli.downloader = manager.NewDownloader(
			cli.cli,
			func(downloader *manager.Downloader) {
				downloader.Concurrency = cli.meta.concurrency
				downloader.Logger = &DefLogger{dummy: !cli.meta.debug}
			},
		)
	}
	{
		_, err := cli.cli.ListBuckets(
			context.Background(),
			&s3.ListBucketsInput{},
		)
		if err != nil {
			return nil, err
		}
	}

	return cli, nil
}

func (s *S3) createBucket(ctx context.Context, bucket string) error {
	_, err := s.cli.CreateBucket(
		ctx,
		&s3.CreateBucketInput{
			Bucket: &bucket,
		},
	)
	return err
}

func (s *S3) upload(ctx context.Context, bucket, objectKey string, data []byte) error {
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

func (s *S3) download(ctx context.Context, bucket, objectKey string) ([]byte, error) {
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
