package aws_s3

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type awsS3Repo struct {
	conn *s3.S3
}

func NewAwsS3Repo(conn *s3.S3) *awsS3Repo {
	return &awsS3Repo{
		conn,
	}
}

func (awsS3 *awsS3Repo) UploadUserProfilePicture(fileName string, file io.ReadSeeker) error {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	if bucketName == "" {
		return errors.New("missing AWS_S3_BUCKET_NAME env variable")
	}

	rootKey := os.Getenv("AWS_S3_ROOT_KEY")

	if rootKey == "" {
		return errors.New("missing AWS_S3_ROOT_KEY env variable")
	}

	filePath := fmt.Sprintf("%v/users/%v", rootKey, fileName)

	_, err := awsS3.conn.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
		Body:   file,
	})

	return err
}

func (awsS3 *awsS3Repo) DeleteUserProfilePicture(fileName string) error {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	if bucketName == "" {
		return errors.New("missing AWS_S3_BUCKET_NAME env variable")
	}

	rootKey := os.Getenv("AWS_S3_ROOT_KEY")

	if rootKey == "" {
		return errors.New("missing AWS_S3_ROOT_KEY env variable")
	}

	filePath := fmt.Sprintf("%v/users/%v", rootKey, fileName)

	_, err := awsS3.conn.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filePath),
	})

	return err

}
