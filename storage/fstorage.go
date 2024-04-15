package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	// PutFile saves the file to the storage and returns the key of the file
	PutFile(file io.Reader, fname, ftype string, fsize int64) (string, error)
}

type s3Storage struct {
	bucketName string
	s3         *S3Client
}

func NewS3Storage(s3 *S3Client, bucketName string) Storage {
	return &s3Storage{
		s3:         s3,
		bucketName: bucketName,
	}
}

func (s *s3Storage) PutFile(file io.Reader, fname, ftype string, fsize int64) (string, error) {
	objectKey := fmt.Sprintf("%s_%d.%s", fname, time.Now().Unix(), ftype)

	_, err := s.s3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(objectKey),
		ContentLength: aws.Int64(fsize),
		ContentType:   aws.String(ftype),
		Body:          file,
	})
	if err != nil {
		return "", err
	}

	return objectKey, nil
}
