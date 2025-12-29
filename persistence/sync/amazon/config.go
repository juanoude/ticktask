package amazon

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const region = "us-east-1"

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
