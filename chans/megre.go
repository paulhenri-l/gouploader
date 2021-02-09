package chans

import (
	"github.com/paulhenri-l/gouploader/contracts"
	"sync"
)

// Merge the output of multiple chans of contracts.UploadResult into one (fan-in)
func MergeChanOfContractsUploadResult(chans ...<-chan contracts.UploadResult) <-chan contracts.UploadResult {
	var wg sync.WaitGroup
	merged := make(chan contracts.UploadResult)

	for _, chanToMerge := range chans {
		go doMergeChanOfContractsUploadResult(chanToMerge, merged, &wg)
		wg.Add(1)
	}

	go waitAndCloseChanOfContractsUploadResult(&wg, merged)

	return merged
}

func doMergeChanOfContractsUploadResult(in <-chan contracts.UploadResult, out chan<- contracts.UploadResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range in {
		out <- item
	}
}

func waitAndCloseChanOfContractsUploadResult(wg *sync.WaitGroup, merged chan contracts.UploadResult) {
	wg.Wait()
	close(merged)
}
