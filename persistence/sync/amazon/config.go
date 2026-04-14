// Package amazon provides AWS S3 operations for database backup and restore.
// It wraps the AWS SDK v2 to provide simple upload/download functionality.
package amazon

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// region is the AWS region for S3 operations.
// Currently hardcoded to us-east-1.
// TODO: Make this configurable via the sync config command.
const region = "us-east-1"

// LoadConfig creates an AWS SDK configuration with the provided credentials.
// Returns the config and a background context for use with AWS API calls.
// Fatally exits if the configuration cannot be loaded.
func LoadConfig(accId string, accSecret string) (aws.Config, context.Context) {
	ctx := context.Background()
	creds := credentials.NewStaticCredentialsProvider(accId, accSecret, "")
	sdkConfig, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		log.Fatal("error during aws config retrieval")
	}

	return sdkConfig, ctx
}
