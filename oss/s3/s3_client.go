package s3

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cheggaaa/pb"
	"io"
	"io/fs"
	"strings"
)

type Client interface {
	BucketFiles(ctx context.Context, bucketName string, prefixHint string) ([]string, error)
	UploadFile(ctx context.Context, bucketName string, remotePath string, localFile fs.File) error
	DownloadFile(ctx context.Context, bucketName string, remotePath string, localFile io.WriterAt) error
	DeleteFile(ctx context.Context, bucketName string, remotePath string) error
}

type S3clientOpt interface {
	apply(*S3client)
}

type s3clientOptFn func(*S3client)

func (f s3clientOptFn) apply(cli *S3client) {
	f(cli)
}

func WithAwsConfig(awsConfig *aws.Config) S3clientOpt {
	return s3clientOptFn(func(cli *S3client) {
		cli.config.awsConfig = awsConfig
	})
}

func WithProgressOutput(progressOutput io.Writer) S3clientOpt {
	return s3clientOptFn(func(cli *S3client) {
		cli.config.progressOutput = progressOutput
	})
}

type S3client struct {
	config struct {
		awsConfig      *aws.Config
		progressOutput io.Writer
	}
	cli *s3.Client
}

func initS3client() *S3client {
	s3m := &S3client{}
	return s3m
}

func NewS3Cli(opts ...S3clientOpt) (*S3client, error) {
	s3c := initS3client()
	{
		for _, opt := range opts {
			opt.apply(s3c)
		}
		switch {
		case s3c.config.awsConfig == nil:
			return nil, errWrapper(IllegalParams)
		case s3c.config.progressOutput == nil:
			return nil, errWrapper(IllegalParams)
		}
	}

	s3c.cli = s3.NewFromConfig(
		*s3c.config.awsConfig,
		func(options *s3.Options) {
			options.UsePathStyle = true
		},
	)

	return s3c, nil
}

func (c *S3client) BucketFiles(ctx context.Context, bucketName string, directoryPrefix string) ([]string, error) {
	if !strings.HasSuffix(directoryPrefix, "/") {
		directoryPrefix = directoryPrefix + "/"
	}
	var (
		continuationToken *string
		truncated         bool
		paths             []string
	)
	for continuationToken, truncated = nil, true; truncated; {
		s3ListChunk, err := c.chunkedBucketList(ctx, bucketName, directoryPrefix, continuationToken)
		if err != nil {
			return []string{}, err
		}
		truncated = s3ListChunk.Truncated
		continuationToken = s3ListChunk.ContinuationToken
		paths = append(paths, s3ListChunk.Paths...)
	}
	return paths, nil
}

type BucketListChunk struct {
	Truncated         bool
	ContinuationToken *string
	CommonPrefixes    []string
	Paths             []string
}

func (c *S3client) chunkedBucketList(ctx context.Context, bucketName string, prefix string, continuationToken *string) (BucketListChunk, error) {
	response, err := c.cli.ListObjectsV2(
		ctx,
		&s3.ListObjectsV2Input{
			Bucket:            aws.String(bucketName),
			ContinuationToken: continuationToken,
			Delimiter:         aws.String("/"),
			Prefix:            aws.String(prefix),
		},
	)
	if err != nil {
		return BucketListChunk{}, err
	}
	commonPrefixes := make([]string, 0, len(response.CommonPrefixes))
	paths := make([]string, 0, len(response.Contents))
	for _, commonPrefix := range response.CommonPrefixes {
		commonPrefixes = append(commonPrefixes, *commonPrefix.Prefix)
	}
	for _, path := range response.Contents {
		paths = append(paths, *path.Key)
	}
	return BucketListChunk{
		Truncated:         response.IsTruncated,
		ContinuationToken: response.NextContinuationToken,
		CommonPrefixes:    commonPrefixes,
		Paths:             paths,
	}, nil
}
func (c *S3client) newProgressBar(total int64) *pb.ProgressBar {
	progress := pb.New64(total)

	progress.Output = c.config.progressOutput
	progress.ShowSpeed = true
	progress.Units = pb.U_BYTES
	progress.NotPrint = true

	return progress.SetWidth(80)
}

func (c *S3client) UploadFile(ctx context.Context, bucketName string, remotePath string, localFile fs.File) error {
	uploader := manager.NewUploader(c.cli)

	stat, err := localFile.Stat()
	if err != nil {
		return err
	}
	progress := c.newProgressBar(stat.Size())

	progress.Start()
	defer progress.Finish()

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(remotePath),
		Body:   progressReader{localFile, progress},
	})
	return err
}

func (c *S3client) DownloadFile(ctx context.Context, bucketName string, remotePath string, localFile io.WriterAt) error {
	object, err := c.cli.HeadObject(
		ctx,
		&s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(remotePath),
		},
	)
	if err != nil {
		return err
	}

	progress := c.newProgressBar(object.ContentLength)

	progress.Start()
	defer progress.Finish()

	safeDown := func() error {
		_, err = manager.NewDownloader(c.cli).Download(
			ctx,
			progressWriterAt{localFile, progress},
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(remotePath),
			},
		)
		return err
	}
	unsafeDown := func() error {
		output, err := c.cli.GetObject(
			ctx,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(remotePath),
			},
		)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(nil)
		if _, err = io.Copy(buf, progressReader{output.Body, progress}); err != nil {
			return err
		}
		defer output.Body.Close()
		_, err = localFile.WriteAt(buf.Bytes(), 0)
		return err
	}
	if unsafe, _ := ctx.Value("unsafe").(bool); unsafe {
		return unsafeDown()
	}
	return safeDown()
}
func (c *S3client) DeleteFile(ctx context.Context, bucketName string, remotePath string) error {
	_, err := c.cli.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(remotePath),
	})

	return err
}
