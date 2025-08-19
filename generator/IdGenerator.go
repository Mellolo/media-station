package generator

import (
	"fmt"
	"github.com/Mellolo/common/cache"
	"github.com/Mellolo/common/errors"
	"strconv"
)

type IdGenerator interface {
	GenerateId(key string) string
}

func NewIdGenerator() *IdGeneratorImpl {
	client, err := cache.GetCache()
	if err != nil {
		panic(errors.WrapError(err, "get generator client failed"))
	}
	return &IdGeneratorImpl{
		client: client,
	}
}

type IdGeneratorImpl struct {
	client cache.Cache
}

func (impl *IdGeneratorImpl) GenerateId(key string) string {
	id, err := impl.client.IncrID(key)
	if err != nil {
		panic(errors.WrapError(err, fmt.Sprintf("generate id with key[%s] failed", key)))
	}
	return fmt.Sprintf("%s_%s", key, strconv.FormatInt(id, 10))
}
