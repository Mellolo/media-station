package cache

import (
	"github.com/mellolo/common/cache"
	"github.com/mellolo/common/errors"
	"time"
)

type DistributedCache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expire time.Duration) error
}

func NewDistributedCache() *DistributedCacheImpl {
	client, err := cache.GetCache()
	if err != nil {
		panic(errors.WrapError(err, "get cache client failed"))
	}
	return &DistributedCacheImpl{
		client: client,
	}
}

type DistributedCacheImpl struct {
	client cache.Cache
}

func (impl DistributedCacheImpl) Get(key string) (interface{}, error) {
	return impl.client.Get(key)
}

func (impl DistributedCacheImpl) Set(key string, value interface{}, expire time.Duration) error {
	return impl.client.Set(key, value, expire)
}
