package storage

import (
	"encoding/json"
	"errors"
	"github.com/alin-io/pkgstore/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
)

type S3Backend struct {
	BaseStorageBackend
	s3Session *session.Session
	s3        *s3.S3

	Bucket string
}

func NewS3Backend() *S3Backend {
	s, err := session.NewSession(&aws.Config{
		Endpoint:         &config.Get().Storage.S3.ApiHost,
		Region:           &config.Get().Storage.S3.Region,
		Credentials:      credentials.NewStaticCredentials(config.Get().Storage.S3.ApiKey, config.Get().Storage.S3.ApiSecret, ""),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}

	return &S3Backend{
		s3Session: s,
		s3:        s3.New(s),

		Bucket: config.Get().Storage.S3.Bucket,
	}
}

func (s *S3Backend) WriteFile(key string, fileMeta interface{}, r io.Reader) error {
	metadata := make(map[string]*string)
	if fileMeta != nil {
		metadataBuffer, err := json.Marshal(fileMeta)
		if err != nil {
			return err
		}
		err = json.Unmarshal(metadataBuffer, &metadata)
		if err != nil {
			return err
		}
	}
	uploader := s3manager.NewUploader(s.s3Session, func(u *s3manager.Uploader) {})
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   r,
	})
	return err
}

func (s *S3Backend) GetFile(key string) (io.ReadCloser, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	obj, err := s.s3.GetObject(getObjectInput)
	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, nil
			case s3.ErrCodeNoSuchKey:
				return nil, nil
			}
		}
		return nil, err
	}
	return obj.Body, nil
}

func (s *S3Backend) CopyFile(fromKey, toKey string) error {
	_, err := s.s3.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s.Bucket),
		CopySource: aws.String(s.Bucket + "/" + fromKey),
		Key:        aws.String(toKey),
	})
	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil
			case s3.ErrCodeNoSuchKey:
				return nil
			}
		}
	}
	return err
}

func (s *S3Backend) DeleteFile(key string) error {
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	_, err := s.s3.DeleteObject(deleteObjectInput)
	if err != nil {
		var aerr awserr.Error
		if errors.As(err, &aerr) {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil
			case s3.ErrCodeNoSuchKey:
				return nil
			}
		}
	}
	return err
}

func (s *S3Backend) GetMetadata(key string, value interface{}) error {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	obj, err := s.s3.GetObject(getObjectInput)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(obj.Body)
	metadataBuffer, err := json.Marshal(obj.Metadata)
	if err != nil {
		return err
	}
	err = json.Unmarshal(metadataBuffer, value)
	return err
}
