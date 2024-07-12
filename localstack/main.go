package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Svc *s3.Client

func init() {
	// Load the SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(fmt.Sprintf("Unable to load SDK config: %v", err))
	}

	// Initialize an S3 client
	s3Svc = s3.NewFromConfig(cfg)
}

func HandleRequest(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {

		getObjectOutput, err := s3Svc.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(record.S3.Bucket.Name),
			Key:    aws.String(record.S3.Object.Key),
		})

		if err != nil {
			fmt.Println(fmt.Errorf("unable to download object: %v, reason: %w", record.S3.Object.Key, err))
			continue
		}

		defer getObjectOutput.Body.Close()

		_, err = s3Svc.PutObject(ctx, &s3.PutObjectInput{
			Bucket:        aws.String("thumbnails"),
			Key:           aws.String(record.S3.Object.Key),
			Body:          getObjectOutput.Body,
			ContentLength: getObjectOutput.ContentLength,
		}, s3.WithAPIOptions(
			v4.SwapComputePayloadSHA256ForUnsignedPayloadMiddleware,
		))

		if err != nil {
			fmt.Println(fmt.Errorf("unable to upload object: %v, reason: %w", record.S3.Object.Key, err))
			continue
		}
	}
}

func main() {
	lambda.Start(HandleRequest)
}
