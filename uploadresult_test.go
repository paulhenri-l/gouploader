package gouploader

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUploadResult_GetFilepath(t *testing.T) {
	fp := "some-file-path"

	ur := UploadResult{}
	ur.Filepath = fp

	assert.Equal(t, fp, ur.GetFilepath())
}

func TestUploadResult_GetSize(t *testing.T) {
	size := int64(1000)

	ur := UploadResult{}
	ur.Size = size

	assert.Equal(t, size, ur.GetSize())
}

func TestUploadResult_GetStart(t *testing.T) {
	start := time.Now()

	ur := UploadResult{}
	ur.Start = start

	assert.Equal(t, start, ur.GetStart())
}

func TestUploadResult_GetEnd(t *testing.T) {
	end := time.Now()

	ur := UploadResult{}
	ur.End = end

	assert.Equal(t, end, ur.GetEnd())
}

func TestUploadResult_GetDuration(t *testing.T) {
	ur := UploadResult{}
	ur.Start = time.Now()
	ur.End = ur.Start.Add(time.Second * 10)

	assert.Equal(t, float64(10), ur.GetDuration().Seconds())
}

func TestUploadResult_GetSpeed(t *testing.T) {
	ur := UploadResult{}
	ur.Start = time.Now()
	ur.End = ur.Start.Add(time.Second * 10)
	ur.Size = (1024 * 1024) * 10

	assert.Equal(t, float64(1), ur.GetSpeed())
}

func TestUploadResult_GetError(t *testing.T) {
	err := errors.New("hello")

	ur := UploadResult{}
	ur.Error = err

	assert.Equal(t, err, ur.GetError())
}

func TestUploadResult_GetError_Nil(t *testing.T) {
	ur := UploadResult{}

	assert.Nil(t, ur.GetError())
}
