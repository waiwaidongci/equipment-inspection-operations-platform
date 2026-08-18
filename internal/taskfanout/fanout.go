package taskfanout

import (
	"sync"
)

func Run(inputs []string, process func(string) error) ([]string, error) {
	n := len(inputs)
	results := make(chan string, n)
	errs := make(chan error, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for _, input := range inputs {
		go func(value string) {
			defer wg.Done()
			if err := process(value); err != nil {
				errs <- err
				return
			}
			results <- value
		}(input)
	}
	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()
	var out []string
	var firstErr error
	for results != nil || errs != nil {
		select {
		case value, ok := <-results:
			if !ok {
				results = nil
				continue
			}
			out = append(out, value)
		case err, ok := <-errs:
			if !ok {
				errs = nil
				continue
			}
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return out, firstErr
}
