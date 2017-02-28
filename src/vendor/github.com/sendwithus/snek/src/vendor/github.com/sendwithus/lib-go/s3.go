package swu

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Service interface {
	Read(bucket string, key string) (io.Reader, error)
	Write(bucket string, key string, reader io.Reader) error
}

type s3ServiceImpl struct {
	S3 *s3.S3
}

func NewS3Service() S3Service {
	return s3ServiceImpl{
		S3: s3.New(getAwsSession()),
	}
}

func (s s3ServiceImpl) Read(bucket string, key string) (io.Reader, error) {

	downloader := s3manager.NewDownloader(getAwsSession())
	buffer := &aws.WriteAtBuffer{}
	_, err := downloader.Download(buffer, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buffer.Bytes()), nil
}

func (s s3ServiceImpl) Write(bucket string, key string, reader io.Reader) error {

	uploader := s3manager.NewUploader(getAwsSession())
	input := s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	}
	_, err := uploader.Upload(&input)

	if err != nil {
		msg := fmt.Sprintf("Failed to upload data to %s/%s: %s", bucket, key, err.Error())
		if multiErr, ok := err.(s3manager.MultiUploadFailure); ok {
			msg = fmt.Sprintf(
				"Failed to upload data to %s/%s, code:%s, msg:%s. uploadId:%s",
				bucket,
				key,
				multiErr.Code(),
				multiErr.Message(),
				multiErr.UploadID())
		}
		return errors.New(msg)
	}

	return nil
}
