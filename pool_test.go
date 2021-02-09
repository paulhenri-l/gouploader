package gouploader

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/gouploader/contracts"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	u1, _ := fakeUploader(t)
	u2, _ := fakeUploader(t)
	m := []contracts.Uploader{u1, u2}
	p := NewPool(m)

	assert.IsType(t, &Pool{}, p)
}

func TestPool_Upload(t *testing.T) {
	u1, _ := fakeUploader(t)
	u2, _ := fakeUploader(t)
	m := []contracts.Uploader{u1, u2}
	p := NewPool(m)

	u1.EXPECT().Upload(gomock.Eq("some-file-1")).Times(1)
	u2.EXPECT().Upload(gomock.Eq("some-file-2")).Times(1)

	_ = p.Upload("some-file-1")
	_ = p.Upload("some-file-2")
}

func TestPool_Upload_Concurrency(t *testing.T) {
	var u1Counter, u2Counter, u3Counter, u4Counter int
	u1, _ := fakeUploader(t)
	u2, _ := fakeUploader(t)
	u3, _ := fakeUploader(t)
	u4, _ := fakeUploader(t)
	m := []contracts.Uploader{u1, u2, u3, u4}
	p := NewPool(m)

	// Update uploaders counter on every call. Since the update is not atomic
	// if the uploader is shared between two goroutines the count won't be
	// correct thus proving the existence of a race condition
	u1.EXPECT().Upload(gomock.Eq("some-file")).AnyTimes().Do(func(_ string) {
		u1Counter = u1Counter + 1
	})

	u2.EXPECT().Upload(gomock.Eq("some-file")).AnyTimes().Do(func(_ string) {
		u2Counter = u2Counter + 1
	})

	u3.EXPECT().Upload(gomock.Eq("some-file")).AnyTimes().Do(func(_ string) {
		u3Counter = u3Counter + 1
	})

	u4.EXPECT().Upload(gomock.Eq("some-file")).AnyTimes().Do(func(_ string) {
		u4Counter = u4Counter + 1
	})

	// Inside 10 goroutines call 100 times upload
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				if err := p.Upload("some-file"); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, 100000, u1Counter+u2Counter+u3Counter+u4Counter)
}

func TestPool_Upload_Broken(t *testing.T) {
	u1, _ := fakeUploader(t)
	m := []contracts.Uploader{u1}
	p := NewPool(m)

	u1.EXPECT().Upload(
		gomock.Eq("some-file-1"),
	).Return(
		&UploadResult{Error: errors.New("some-error")},
	)

	res := p.Upload("some-file-1")

	assert.EqualError(t, errors.Cause(res.GetError()), "some-error")
}

func TestPool_Upload_Wait(t *testing.T) {
	u1, _ := fakeUploader(t)
	u2, _ := fakeUploader(t)
	m := []contracts.Uploader{u1, u2}
	p := NewPool(m)
	wait := make(chan struct{})

	go func() {
		u1, _ := p.take()
		u2, _ := p.take()
		wait <- struct{}{}
		time.Sleep(200 * time.Millisecond)
		p.put(u1)
		p.put(u2)
	}()

	u1.EXPECT().Upload(gomock.Eq("some-file")).Times(1)

	// Wait for pool to be empty then attempt to write
	// this will block until all managers are back in the pool
	<-wait
	_ = p.Upload("some-file")
}

func TestPool_Close(t *testing.T) {
	u1, _ := fakeUploader(t)
	m := []contracts.Uploader{u1}
	p := NewPool(m)

	p.Close()
	res := p.Upload("some-file-1")

	assert.EqualError(t, errors.Cause(res.GetError()), "pool closed")
}

func TestPool_Close_Wait(t *testing.T) {
	uc := make(chan contracts.Uploader, 2)
	u1, _ := fakeUploader(t)
	u2, _ := fakeUploader(t)
	m := []contracts.Uploader{u1, u2}
	p := NewPool(m)
	wait := make(chan struct{})

	go func() {
		// Empty pool
		u1, _ := p.take()
		u2, _ := p.take()

		// Make upload
		wait <- struct{}{}

		// Put back in 50ms
		uc <- u1
		uc <- u2
		close(uc)

		// Close
		p.Close()
	}()

	go func() {
		for u := range uc {
			time.Sleep(50 * time.Millisecond)
			p.put(u)
		}
	}()

	u1.EXPECT().Upload(gomock.Eq("some-file")).Times(0)

	// Wait for pool to be empty then attempt to upload one file. The pool will
	// be closed right before the uploaders are put back in the pool but after
	// the call to upload
	<-wait
	res := p.Upload("some-file")
	assert.EqualError(t, errors.Cause(res.GetError()), "pool closed")
}

func TestPool_CloseWithContext(t *testing.T) {
	u1, _ := fakeUploader(t)
	u2, _ := fakeUploader(t)
	m := []contracts.Uploader{u1, u2}
	p := NewPool(m)

	// Take the uploaders out of the pool
	t1, _ := p.take()
	t2, _ := p.take()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	go func() {
		<-time.NewTimer(time.Millisecond * 200).C
		p.put(t1)
		p.put(t2)
	}()

	err := p.CloseWithContext(ctx)

	assert.Nil(t, err)
	cancel()
}

func TestPool_CloseWithContext_Timeout(t *testing.T) {
	u1, _ := fakeUploader(t)
	u2, _ := fakeUploader(t)
	m := []contracts.Uploader{u1, u2}
	p := NewPool(m)

	// Take the uploaders out of the pool
	_, _ = p.take()
	_, _ = p.take()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	err := p.CloseWithContext(ctx)

	assert.EqualError(t, err, "context deadline exceeded")
	cancel()
}
