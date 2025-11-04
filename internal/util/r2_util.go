package util

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func NewR2Client(cfg *viper.Viper) (*s3.Client, error) {
	accessKey := cfg.GetString("cloudflare_r2.access_key")
	secretKey := cfg.GetString("cloudflare_r2.secret_key")
	endpoint := cfg.GetString("cloudflare_r2.endpoint")

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: "auto",
			}, nil
		},
	)

	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(awsCfg), nil
}
