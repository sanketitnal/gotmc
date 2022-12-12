# Golang TTL Memoization in-memory Cache
Golang TTL Memoization in-memory Cache

## How to use ?
```
package main

import (
    "log"
    "time"
    "github.com/sanketitnal/gotmc/tmc"
)

func main() {
    /*
    * -> type Func func(key string) (interface{}, error)
    *
    * -> tmc.NewTMCache(fun Func, cleanupTimeout time.Duration) *TMCache
    *  
    */
    cacheCleanupDeadline := time.Second*30
    cache := tmc.NewTMCache(tmc.HttpGetBody, cacheCleanupDeadline)
    
    keyTTL := time.Second*30

    // value, cacheHit?, error?
    val, chit, err := cache.Get("http://www.google.com", keyTTL)
    /*
    * If val doesn't exists in cache, fun will be called with key and result is stored in cache
    */

    if err != nil {
        log.Error(err)
    }
    
    // value, cacheHit?, error?
    val, chit, err = cache.Get("http://www.google.com", keyTTL)
    // This should return chit = true
    
    fmt.Println(val, chit, err)
    
    cache.Close()
}
```

## Methods
- ```NewTMCache(fun Func, cleanupTimeout time.Duration) *TMCache```
Returns instance of cache. Cleanup will be run after every cleanupTimeout duration.
- ```Get(key string, ttl time.Duration) (interface{}, bool, error)```
Returns value cached for key. If value isn't cached, value is obtained by calling fun(key).
ttl is duration after which this key will be cleaned up.
- ```Del(key string)``` Deletes key.
- ```EraseAll()``` Deletes all keys in the cache.
- ```Close()``` Deletes all keys in cache, and makes it nil.

## Benchmark
- Benchmark was run with fun = tmc.HttpGetBody(url string) function. 1000 urls were requested 100 times with and without cache. Results:
```
cpu: Intel(R) Core(TM) i5-7300HQ CPU @ 2.50GHz
BenchmarkCache-4        1000000000               0.03700 ns/op
BenchmarkNoCache-4      1000000000               0.1152 ns/op
PASS
ok      github.com/sanketitnal/gotmc/tmc        2.926s
```
- Note that benchmark results may vary for every run and from system to system.
