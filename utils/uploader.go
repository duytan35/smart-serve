package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Uploader struct {
	Client     *s3.Client
	BucketName string
}

type File struct {
	ID       uuid.UUID
	Name     string
	MineType string
}

func NewS3Uploader() (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("S3_BUCKET_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), "")),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		Client:     client,
		BucketName: os.Getenv("S3_BUCKET_NAME"),
	}, nil
}

func (u *S3Uploader) UploadFile(file multipart.File, fileHeader *multipart.FileHeader, key uuid.UUID) (File, error) {
	defer file.Close()

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	contentType := http.DetectContentType(buffer)

	_, err := u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(u.BucketName),
		Key:           aws.String(key.String()),
		Body:          bytes.NewReader(buffer),
		ContentLength: &size,
		ContentType:   aws.String(contentType),
	})

	if err != nil {
		return File{}, fmt.Errorf("failed to upload file to S3, %v", err)
	}

	return File{
		ID:       key,
		Name:     fileHeader.Filename,
		MineType: http.DetectContentType(buffer),
	}, nil
}

func (u *S3Uploader) RemoveFile(key uuid.UUID) error {
	_, err := u.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(key.String()),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3, %v", err)
	}

	return nil
}

var Uploader *S3Uploader

func InitS3Uploader() {
	var err error
	Uploader, err = NewS3Uploader()
	if err != nil {
		log.Fatalf("Failed to create S3 uploader: %v", err)
	}
}
