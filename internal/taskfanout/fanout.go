package taskfanout

import (
	"errors"
	"sync"
)

func Run(inputs []string, process func(string) error) ([]string, error) {
	results := make(chan string)
	errorsCh := make(chan error)
	var wg sync.WaitGroup
	for _, input := range inputs {
		go func(value string) {
			wg.Add(1)
			defer wg.Done()
			if err := process(value); err != nil {
				errorsCh <- err
				return
			}
			results <- value
		}(input)
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(results); close(done) }()
	var out []string
	for value := range results {
		out = append(out, value)
	}
	select {
	case <-done:
		return out, nil
	default:
		return out, errors.New("fanout incomplete")
	}
}
