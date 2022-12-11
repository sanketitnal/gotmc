package tmc

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkCache(b *testing.B) {
	var wg sync.WaitGroup
	hour := time.Second * 60 * 60
	cache := NewTMCache(HttpGetBody, hour)
	times := 100
	for times != 0 {
		for _, url := range TestUrls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				cache.Get(url, hour)
			}(url)
		}
		times--
	}
	wg.Wait()
}
