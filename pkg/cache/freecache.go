package cache

import (
	"time"

	"github.com/coocood/freecache"
)

type freeCache struct {
	cache *freecache.Cache
}

func (f *freeCache) Set(key string, value []byte, ttl time.Duration) error {
	return f.cache.Set([]byte(key), value, int(ttl.Seconds()))
}

func (f *freeCache) Get(key string) ([]byte, error) {
	val, err := f.cache.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (f *freeCache) Del(key string) error {
	f.cache.Del([]byte(key))
	return nil
}

func (f *freeCache) Has(key string) bool {
	_, err := f.cache.Get([]byte(key))
	return err == nil
}
