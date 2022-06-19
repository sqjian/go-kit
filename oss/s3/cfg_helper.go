package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type s3configOpt interface {
	apply(*s3Config)
}

type s3configOptFn func(*s3Config)

func (f s3configOptFn) apply(cfg *s3Config) {
	f(cfg)
}
func WithAccessKey(accessKey string) s3configOpt {
	return s3configOptFn(func(cfg *s3Config) {
		cfg.accessKey = accessKey
	})
}
func WithSecretKey(secretKey string) s3configOpt {
	return s3configOptFn(func(cfg *s3Config) {
		cfg.secretKey = secretKey
	})
}
func WithEndpoint(endpoint string) s3configOpt {
	return s3configOptFn(func(cfg *s3Config) {
		cfg.endpoint = endpoint
	})
}

type s3Config struct {
	accessKey string
	secretKey string
	endpoint  string
}

func initS3config() *s3Config {
	cfg := &s3Config{}
	return cfg
}
func NewAwsConfig(opts ...s3configOpt) (*aws.Config, error) {
	s3cfg := initS3config()
	{
		for _, opt := range opts {
			opt.apply(s3cfg)
		}
		switch {
		case s3cfg.accessKey == "" || s3cfg.secretKey == "":
			return nil, errWrapper(IllegalParams)
		case s3cfg.endpoint == "":
			return nil, errWrapper(IllegalParams)
		}
	}

	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s3cfg.accessKey,
				s3cfg.secretKey,
				"",
			)),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: s3cfg.endpoint,
					}, nil
				})),
	)
	if err != nil {
		return nil, err
	}
	return &awsCfg, nil
}
