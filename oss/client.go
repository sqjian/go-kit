package actions

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func genCli() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				"accessKey",
				"secretKey",
				"",
			)),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...any) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL: "endpoint",
					}, nil
				})),
	)
	checkErr(err)

	_ = s3.NewFromConfig(
		awsCfg,
		func(options *s3.Options) {
			options.UsePathStyle = true
		},
	)
}
