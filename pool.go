package gouploader

import (
	"context"
	"github.com/paulhenri-l/gouploader/contracts"
	"github.com/pkg/errors"
	"sync"
)

type Pool struct {
	mtx       *sync.Mutex
	cnd       *sync.Cond
	size      int
	uploaders []contracts.Uploader
	closed    bool
}

func NewPool(uploaders []contracts.Uploader) *Pool {
	mtx := &sync.Mutex{}
	cnd := sync.NewCond(mtx)

	return &Pool{
		mtx:       mtx,
		cnd:       cnd,
		size:      len(uploaders),
		uploaders: uploaders,
	}
}

func (p *Pool) Upload(filepath string) contracts.UploadResult {
	u, err := p.take()
	if err != nil {
		return &UploadResult{Filepath: filepath, Error: err}
	}

	res := u.Upload(filepath)
	p.put(u)

	return res
}

func (p *Pool) Close() {
	_ = p.CloseWithContext(context.Background())
}

func (p *Pool) CloseWithContext(ctx context.Context) error {
	// Make it possible to kill the goroutine if we go over the ctx deadline
	var killGoRoutine bool

	p.mtx.Lock()
	p.closed = true
	p.mtx.Unlock()
	closed := make(chan struct{}, 1)

	// Wait for uploader to be back in the pool. When they are back send to the
	// closed chan and this will unblock the calling function.
	go func() {
		p.mtx.Lock()
		defer p.mtx.Unlock()

		for len(p.uploaders) < p.size && !killGoRoutine {
			p.cnd.Wait()
		}

		closed <- struct{}{}
	}()

	// Wait for either ctx deadline or closed pool
	select {
	case <-ctx.Done():
		p.mtx.Lock()
		killGoRoutine = true
		p.mtx.Unlock()
		return ctx.Err()
	case <-closed:
		return nil
	}
}

func (p *Pool) put(u contracts.Uploader) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	defer p.cnd.Signal()

	p.uploaders = append(p.uploaders, u)
}

func (p *Pool) take() (contracts.Uploader, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	for len(p.uploaders) <= 0 {
		p.cnd.Wait()
	}

	if p.closed {
		return nil, errors.New("pool closed")
	}

	u, uploaders := p.uploaders[0], p.uploaders[1:]
	p.uploaders = uploaders

	return u, nil
}
