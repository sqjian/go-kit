package oss

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/sqjian/go-kit/log"
	"io"
	"os"
)

func NewS3Stub(accessKey, secretKey, endpoint string) (*S3Stub, error) {
	configInst, configInstErr := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			)),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...any) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint}, nil
				})),
	)
	if configInstErr != nil {
		return nil, configInstErr
	}
	return &S3Stub{s3Client: s3.NewFromConfig(configInst, func(options *s3.Options) { options.UsePathStyle = true })}, nil

}

type S3Stub struct {
	s3Client *s3.Client
	logger   log.Log
}

// ListBuckets lists the buckets in the current account.
func (s *S3Stub) ListBuckets() ([]types.Bucket, error) {
	result, err := s.s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	var buckets []types.Bucket
	if err != nil {
		s.logger.Errorf("Couldn't list buckets for your account. Here's why: %v", err)
	} else {
		buckets = result.Buckets
	}
	return buckets, err
}

// BucketExists checks whether a bucket exists in the current account.
func (s *S3Stub) BucketExists(bucketName string) (bool, error) {
	_, err := s.s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			var notFound *types.NotFound
			switch {
			case errors.As(apiError, &notFound):
				s.logger.Errorf("Bucket %v is available.", bucketName)
				exists = false
				err = nil
			default:
				s.logger.Infof("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v", bucketName, err)
			}
		}
	} else {
		s.logger.Infof("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}

// CreateBucket creates a bucket with the specified name in the specified Region.
func (s *S3Stub) CreateBucket(name string, region string) error {
	_, err := s.s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		s.logger.Errorf("Couldn't create bucket %v in Region %v. Here's why: %v",
			name, region, err)
	}
	return err
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (s *S3Stub) UploadFile(bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		s.logger.Errorf("Couldn't open file %v to upload. Here's why: %v", fileName, err)
	} else {
		defer file.Close()
		_, err := s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			s.logger.Errorf("Couldn't upload file %v to %v:%v. Here's why: %v",
				fileName, bucketName, objectKey, err)
		}
	}
	return err
}

// UploadLargeObject uses an upload manager to upload data to an object in a bucket.
// The upload manager breaks large data into parts and uploads the parts concurrently.
func (s *S3Stub) UploadLargeObject(bucketName string, objectKey string, largeObject []byte) error {
	largeBuffer := bytes.NewReader(largeObject)
	var partMiBs int64 = 10
	uploader := manager.NewUploader(s.s3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   largeBuffer,
	})
	if err != nil {
		s.logger.Errorf("Couldn't upload large object to %v:%v. Here's why: %v",
			bucketName, objectKey, err)
	}

	return err
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (s *S3Stub) DownloadFile(bucketName string, objectKey string, fileName string) error {
	result, err := s.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		s.logger.Errorf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		s.logger.Errorf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		s.logger.Errorf("Couldn't read object body from %v. Here's why: %v", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

// DownloadLargeObject uses a download manager to download an object from a bucket.
// The download manager gets the data in parts and writes them to a buffer until all
// the data has been downloaded.
func (s *S3Stub) DownloadLargeObject(bucketName string, objectKey string) ([]byte, error) {
	var partMiBs int64 = 10
	downloader := manager.NewDownloader(s.s3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})
	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		s.logger.Errorf("Couldn't download large object from %v:%v. Here's why: %v",
			bucketName, objectKey, err)
	}
	return buffer.Bytes(), err
}

// CopyToFolder copies an object in a bucket to a subfolder in the same bucket.
func (s *S3Stub) CopyToFolder(bucketName string, objectKey string, folderName string) error {
	_, err := s.s3Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(fmt.Sprintf("%v/%v", bucketName, objectKey)),
		Key:        aws.String(fmt.Sprintf("%v/%v", folderName, objectKey)),
	})
	if err != nil {
		s.logger.Errorf("Couldn't copy object from %v:%v to %v:%v/%v. Here's why: %v",
			bucketName, objectKey, bucketName, folderName, objectKey, err)
	}
	return err
}

// ListObjects lists the objects in a bucket.
func (s *S3Stub) ListObjects(bucketName string) ([]types.Object, error) {
	result, err := s.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	var contents []types.Object
	if err != nil {
		s.logger.Errorf("Couldn't list objects in bucket %v. Here's why: %v", bucketName, err)
	} else {
		contents = result.Contents
	}
	return contents, err
}

// DeleteObjects deletes a list of objects from a bucket.
func (s *S3Stub) DeleteObjects(bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := s.s3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})
	if err != nil {
		s.logger.Errorf("Couldn't delete objects from bucket %v. Here's why: %v", bucketName, err)
	}
	return err
}

// DeleteBucket deletes a bucket. The bucket must be empty or an error is returned.
func (s *S3Stub) DeleteBucket(bucketName string) error {
	_, err := s.s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName)})
	if err != nil {
		s.logger.Errorf("Couldn't delete bucket %v. Here's why: %v", bucketName, err)
	}
	return err
}
