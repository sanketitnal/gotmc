package tmc

import (
	"testing"
	"time"
)

const (
	second = time.Second
	minute = 60 * second
	hour   = 60 * minute
)

func TestKeyIsCached(t *testing.T) {
	cleanupTimeout := 10 * second
	cache := NewTMCache(func(key string) (interface{}, error) {
		return "value", nil
	}, cleanupTimeout)

	_, _, err := cache.Get("key", cleanupTimeout)
	val, chit, erro := cache.Get("key", cleanupTimeout)
	if err != nil || erro != nil {
		t.Error("Error: cache.Get('key', cleanupTimeout=10seconds)")
	} else if chit != true {
		t.Error("Error: cache hit = false; expected true")
	} else if val != "value" {
		t.Error("Error: expected val = 'value'")
	}
}

func TestKeyCleanedAfterTTL(t *testing.T) {
	cleanupTimeout := 1 * second

	cache := NewTMCache(func(key string) (interface{}, error) {
		return "value", nil
	}, cleanupTimeout)

	_, _, err := cache.Get("key", cleanupTimeout)

	time.Sleep(1 * time.Second)

	val, chit, erro := cache.Get("key", cleanupTimeout)

	if err != nil || erro != nil {
		t.Error("Error: cache.Get('key', cleanupTimeout=10seconds)")
	} else if chit != false {
		t.Error("Error: cache hit = true; expected false")
	} else if val != "value" {
		t.Error("Error: expected val = 'value'")
	}

}
