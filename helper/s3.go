package helper

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type S3Helper struct {
	client *s3.S3
	bucket string
	region string
}

var s3HelperInstance *S3Helper

func newS3Helper(accessKey, secretKey, region, bucket string) *S3Helper {
	session := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))
	return &S3Helper{
		client: s3.New(session),
		region: region,
		bucket: bucket,
	}
}

func GetS3Instance() *S3Helper {
	if s3HelperInstance == nil {
		v := viper.GetViper()
		s3HelperInstance = newS3Helper(
			v.GetString("S3_ACCESS_KEY"),
			v.GetString("S3_SECRET_KEY"),
			v.GetString("S3_REGION"),
			v.GetString("S3_BUCKET"),
		)
	}
	return s3HelperInstance
}

func (s *S3Helper) CreateFolder(folderName string, isPublic bool) error {
	config := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(folderName + "/"),
	}
	if isPublic {
		config.ACL = aws.String("public-read")
	}
	_, err := s.client.PutObject(config)
	return err
}

func (s *S3Helper) DeleteFolder(folderName string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(folderName + "/"),
	})
	return err
}

func (s *S3Helper) UploadFile(folderName, fileName string, fileBytes []byte, isPublic bool) (string, error) {
	obj := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(folderName + "/" + fileName),
		Body:   bytes.NewReader(fileBytes),
	}
	if isPublic {
		obj.ACL = aws.String("public-read")
	}
	_, err := s.client.PutObject(obj)
	url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s/%s", s.bucket, s.region, folderName, fileName)
	return url, err
}
