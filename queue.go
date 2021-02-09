package gouploader

import (
	"context"
	"github.com/paulhenri-l/gouploader/chans"
	"github.com/paulhenri-l/gouploader/contracts"
)

type Queue struct {
	parallelism int
	uploader    contracts.Uploader
}

func NewQueue(u contracts.Uploader, parallelism int) *Queue {
	return &Queue{
		uploader:    u,
		parallelism: parallelism,
	}
}

func (q *Queue) Start(ctx context.Context, fileQueue <-chan string) <-chan contracts.UploadResult {
	var outs []<-chan contracts.UploadResult

	if q.parallelism == 0 {
		q.parallelism = 1
	}

	for i := 0; i < q.parallelism; i++ {
		out := make(chan contracts.UploadResult)
		go q.doUpload(ctx, fileQueue, out)
		outs = append(outs, out)
	}

	return chans.MergeChanOfContractsUploadResult(outs...)
}

func (q *Queue) doUpload(
	ctx context.Context,
	fileQueue <-chan string,
	out chan contracts.UploadResult,
) {
	defer close(out)

	for f := range chans.OrDoneString(ctx, fileQueue) {
		out <- q.uploader.Upload(f)
	}
}
