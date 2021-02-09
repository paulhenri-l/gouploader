package chans

import "context"

// Read from the given channel of string as long as the ctx is not done.
func OrDoneString(ctx context.Context, input <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-input:
				if ok != true {
					return
				}

				out <- item
			}
		}
	}()

	return out
}
