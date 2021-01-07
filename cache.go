package bigcache

import (
	"time"

	"github.com/allegro/bigcache/v2"
)

type Cache struct {
	cache *bigcache.BigCache
}

func NewCache(config bigcache.Config) (*Cache, error) {
	cache, err := bigcache.NewBigCache(config)
	return &Cache{cache: cache}, err
}

func (c *Cache) Set(key string, value byte, TTL time.Duration) error {
	timeBinary, err := time.Now().Add(TTL).MarshalBinary()
	if err != nil {
		return err
	}
	v := append(timeBinary, value)
	return c.cache.Set(key, v)
}

func (c *Cache) Get(key string) ([]byte, error) {
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	var bestBefore time.Time
	if err = bestBefore.UnmarshalBinary(value[:15]); err != nil {
		return nil, err
	}

	if time.Now().After(bestBefore) {
		return nil, bigcache.ErrEntryNotFound
	}
	return value[15:], nil
}
