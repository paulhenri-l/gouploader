package gouploader

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/paulhenri-l/gouploader/contracts"
	"github.com/pkg/errors"
	"os"
	"path"
	"time"
)

type S3 struct {
	bucket     string
	s3Uploader s3manageriface.UploaderAPI
}

func NewS3(bucket string, s3Uploader s3manageriface.UploaderAPI) *S3 {
	return &S3{
		bucket:     bucket,
		s3Uploader: s3Uploader,
	}
}

func (u *S3) Upload(file string) contracts.UploadResult {
	ur := &UploadResult{Filepath: file}

	fi, err := os.Stat(file)
	if err != nil {
		ur.Error = errors.Wrap(err, "cannot get file stat")
		return ur
	}

	ur.Size = fi.Size()

	f, err := os.Open(file)
	if err != nil {
		ur.Error = errors.Wrap(err, "cannot open file")
		return ur
	}

	ur.Start = time.Now()
	_, err = u.s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(path.Base(file)),
		Body:   f,
	})
	ur.End = time.Now()

	if err != nil {
		ur.Error = errors.Wrap(err, "S3 error")
		return ur
	}

	return ur
}
