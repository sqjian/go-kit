package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type OptionFunc func(*s3Config)

func WithAccessKey(accessKey string) OptionFunc {
	return func(cfg *s3Config) {
		cfg.accessKey = accessKey
	}
}
func WithSecretKey(secretKey string) OptionFunc {
	return func(cfg *s3Config) {
		cfg.secretKey = secretKey
	}
}
func WithEndpoint(endpoint string) OptionFunc {
	return func(cfg *s3Config) {
		cfg.endpoint = endpoint
	}
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
func NewAwsConfig(opts ...OptionFunc) (*aws.Config, error) {
	s3cfg := initS3config()
	{
		for _, opt := range opts {
			opt(s3cfg)
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
				func(service, region string, options ...any) (aws.Endpoint, error) {
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
