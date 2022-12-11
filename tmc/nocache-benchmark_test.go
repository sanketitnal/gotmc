package tmc

import (
	"sync"
	"testing"
)

func BenchmarkNoCache(b *testing.B) {
	var wg sync.WaitGroup
	times := 100

	for times != 0 {
		for _, url := range TestUrls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				HttpGetBody(url)
			}(url)
		}
		times--
	}

	wg.Wait()
}
