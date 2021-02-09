package gouploader

import (
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/gouploader/mocks/contracts"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestS3_Upload(t *testing.T) {
	m, _ := fakeS3Uploader(t)
	f := fakeFile(t, "hello")
	uploader := NewS3("some-bucket", m)
	matcher := newUploadMatcher(f, "some-bucket")

	m.EXPECT().Upload(gomock.All(matcher)).Return(&s3manager.UploadOutput{}, nil)
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

func fakeS3Uploader(t *testing.T) (*contracts.MockUploaderAPI, *gomock.Controller) {
	ctl := gomock.NewController(t)
	t.Cleanup(func() {
		ctl.Finish()
	})

	return contracts.NewMockUploaderAPI(ctl), ctl
}

type uploadMatcher struct {
	filePath string
	bucket   string
}

func newUploadMatcher(filePath, bucket string) *uploadMatcher {
	return &uploadMatcher{
		filePath: filePath,
		bucket:   bucket,
	}
}

func (u uploadMatcher) Matches(x interface{}) bool {
	ui, ok := x.(*s3manager.UploadInput)
	if !ok {
		return false
	}

	expectedKey := path.Base(u.filePath)
	if *ui.Key != expectedKey {
		return false
	}

	if *ui.Bucket != u.bucket {
		return false
	}

	b, err := ioutil.ReadAll(ui.Body)
	if err != nil {
		return false
	}

	f, err := os.Open(u.filePath)
	if err != nil {
		return false
	}

	expectedContents, err := ioutil.ReadAll(f)
	if err != nil {
		return false
	}

	return string(b) == string(expectedContents)
}

func (u uploadMatcher) String() string {
	return "Upload input did not match with uploaded file"
}
