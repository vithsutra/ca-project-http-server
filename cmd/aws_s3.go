package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3Connection struct {
	s3Client *s3.S3
}

func NewS3Connection() *s3Connection {
	awsRegion := os.Getenv("AWS_S3_REGION")

	if awsRegion == "" {
		log.Fatalln("missing AWS_S3_REGION env variable")
	}

	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")

	if awsAccessKeyId == "" {
		log.Fatalln("missing AWS_ACCESS_KEY_ID env variable")
	}

	awsSecretAccessKey := os.Getenv("AWS_SECRETE_ACCESS_KEY")

	if awsSecretAccessKey == "" {
		log.Fatalln("missing AWS_SECRETE_ACCESS_KEY env variable")
	}

	awsS3BucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	if awsS3BucketName == "" {
		log.Fatalln("missing AWS_S3_BUCKET_NAME env variable")
	}

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(
				awsAccessKeyId,
				awsSecretAccessKey,
				"",
			),
		},
	)

	if err != nil {
		log.Fatalln("error occurred with s3, Error: ", err.Error())
	}

	return &s3Connection{
		s3Client: s3.New(sess),
	}
}
