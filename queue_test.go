package gouploader

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/gouploader/contracts"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewQueue(t *testing.T) {
	m, _ := fakeUploader(t)
	q := NewQueue(m, 1)

	assert.IsType(t, q, &Queue{})
}

func TestQueue_Start_AndClose(t *testing.T) {
	m, _ := fakeUploader(t)
	q := NewQueue(m, 1)
	fq := make(chan string)

	o := q.Start(context.Background(), fq)

	assert.IsType(t, o, make(<-chan contracts.UploadResult))
	close(fq)
	<-o
}

func TestQueue_Start_StopsWhenContextCanceled(t *testing.T) {
	m, _ := fakeUploader(t)
	q := NewQueue(m, 1)
	fq := make(chan string)

	ctx, cancel := context.WithCancel(context.Background())
	o := q.Start(ctx, fq)
	cancel()
	<-o
}

func TestQueue_Upload(t *testing.T) {
	m, _ := fakeUploader(t)
	q := NewQueue(m, 1)
	fq := make(chan string)

	o := q.Start(context.Background(), fq)
	m.EXPECT().Upload(gomock.Eq("some-file")).Times(1).Return(&UploadResult{})

	fq <- "some-file"
	<-o

	close(fq)
	<-o
}

func TestQueue_Upload_Error(t *testing.T) {
	m, _ := fakeUploader(t)
	q := NewQueue(m, 1)
	fq := make(chan string)

	o := q.Start(context.Background(), fq)
	m.EXPECT().Upload(
		gomock.Eq("some-file"),
	).Return(
		&UploadResult{Error: errors.New("bad-things-happened")},
	)

	fq <- "some-file"
	res := <-o

	assert.EqualError(t, res.GetError(), "bad-things-happened")
	close(fq)
}
