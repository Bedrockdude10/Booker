// cache/cache.go
package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

// Global cache instance
var Cache *ristretto.Cache

// Initialize cache once
func Init() error {
	var err error
	Cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 28,
		BufferItems: 64,
	})
	return err
}

// Simple cache functions
func Get(key string) (interface{}, bool) {
	return Cache.Get(key)
}

func Set(key string, value interface{}, ttl time.Duration) {
	Cache.SetWithTTL(key, value, 1, ttl)
}

func Del(key string) {
	Cache.Del(key)
}
