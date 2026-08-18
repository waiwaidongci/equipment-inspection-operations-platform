package taskfanout

import "sync"

func Run(inputs []string, process func(string) error) ([]string, error) {
	results := make(chan string, len(inputs))
	errorsCh := make(chan error, len(inputs))
	var wg sync.WaitGroup
	for _, input := range inputs {
		wg.Add(1)
		go func(value string) {
			defer wg.Done()
			if err := process(value); err != nil {
				errorsCh <- err
				return
			}
			results <- value
		}(input)
	}
	wg.Wait()
	close(results)
	close(errorsCh)
	var out []string
	for value := range results {
		out = append(out, value)
	}
	for err := range errorsCh {
		if err != nil {
			return out, err
		}
	}
	return out, nil
}
