package cache

import (
	"errors"
	"time"

	"github.com/coocood/freecache"

	"github.com/algorithm9/flash-deal/internal/model"
)

type Cache interface {
	Set(key string, value []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
	Del(key string) error
	Has(key string) bool
}

func NewCache(cfg *model.Cache) Cache {
	return &freeCache{
		cache: freecache.NewCache(cfg.MaxBytes),
	}
}

func IsNotFound(err error) bool {
	if errors.Is(err, freecache.ErrNotFound) {
		return true
	}
	return false
}
