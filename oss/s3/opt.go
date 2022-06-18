package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"io"
)

type Option interface {
	apply(*s3client)
}

type optionFunc func(*s3client)

func (f optionFunc) apply(s3c *s3client) {
	f(s3c)
}

func WithAwsConfig(awsConfig *aws.Config) Option {
	return optionFunc(func(cli *s3client) {
		cli.meta.awsConfig = awsConfig
	})
}

func WithProgressOutput(progressOutput io.Writer) Option {
	return optionFunc(func(cli *s3client) {
		cli.meta.progressOutput = progressOutput
	})
}
