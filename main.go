package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sanketitnal/gotmc/tmc"
)

const (
	second = time.Second
	minute = 60 * second
	hour   = 60 * minute
)

func main() {
	cache := tmc.NewTMCache(tmc.HttpGetBody, hour)
	var wg sync.WaitGroup
	for _, url := range tmc.TestUrls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			start := time.Now()
			_, chit, err := cache.Get(url, 100000)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s ___ cache_hit?: %t\n", url, time.Since(start), chit)
		}(url)
	}

	for _, url := range tmc.TestUrls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			start := time.Now()
			_, chit, err := cache.Get(url, 100000)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s ___ cache_hit?: %t\n", url, time.Since(start), chit)
		}(url)
	}
	wg.Wait()

}
