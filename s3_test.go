package gouploader

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/gouploader/mocks"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"path"
	"testing"
)

func TestS3_Upload(t *testing.T) {
	m, _ := fakeS3Uploader(t)
	f := fakeFile(t, "hello")
	uploader := NewS3("some-bucket", m)

	m.EXPECT().Upload(gomock.Any()).DoAndReturn(func(input *s3manager.UploadInput) (*s3manager.UploadOutput, error) {
		assert.Equal(t, "some-bucket", *input.Bucket)
		assert.Equal(t, path.Base(f), *input.Key)

		b, _ := io.ReadAll(input.Body)
		assert.Equal(t, "hello", string(b))

		return &s3manager.UploadOutput{}, nil
	})

	res := uploader.Upload(f)

	assert.NoError(t, res.GetError())
}

func TestS3_Upload_BadFile(t *testing.T) {
	m, _ := fakeS3Uploader(t)
	uploader := NewS3("some-bucket", m)

	res := uploader.Upload("i-dont-exists")

	assert.Error(t, res.GetError())
	assert.Contains(t, res.GetError().Error(), "cannot get file stat")
}

func TestS3_Upload_S3Error(t *testing.T) {
	m, _ := fakeS3Uploader(t)
	f := fakeFile(t, "hello")
	uploader := NewS3("some-bucket", m)
	m.EXPECT().Upload(gomock.Any()).Return(
		&s3manager.UploadOutput{}, errors.New("not good"),
	)

	res := uploader.Upload(f)

	assert.Error(t, res.GetError())
	assert.Contains(t, res.GetError().Error(), "S3 error")
}

func fakeS3Uploader(t *testing.T) (*mocks.MockUploaderAPI, *gomock.Controller) {
	ctl := gomock.NewController(t)
	t.Cleanup(func() {
		ctl.Finish()
	})

	return mocks.NewMockUploaderAPI(ctl), ctl
}

func fakeFile(t *testing.T, contents string) string {
	tmp := t.TempDir()
	guid := xid.New()
	fp := fmt.Sprintf("%s/%s_%s", tmp, "fake_file", guid.String())

	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}

	_, err = f.Write([]byte(contents))
	if err != nil {
		panic(err)
	}

	return fp
}
