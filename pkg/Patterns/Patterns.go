package patterns

import "sync"

func MergeChannels[T any](channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	outputCh := make(chan T)
	for _, channel := range channels {
		go func() {
			defer wg.Done()
			for value := range channel {
				outputCh <- value
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func SplitChannel[T any](inputCh <-chan T, n int) []<-chan T {
	outputChs := make([]chan T, n)
	for i := range n {
		outputChs[i] = make(chan T)
	}

	go func() {
		idx := 0
		for value := range inputCh {
			outputChs[idx] <- value
			idx = (idx + 1) % n
		}

		for _, ch := range outputChs {
			close(ch)
		}
	}()

	resultChs := make([]<-chan T, n)
	for i := range n {
		resultChs[i] = outputChs[i]
	}

	return resultChs
}
