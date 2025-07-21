package initialize

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
	"github.com/mellolo/common/cache"
)

func InitCache() {
	client := redis.NewClient(&redis.Options{
		Addr:     web.AppConfig.DefaultString("redis::redis_end_point", "localhost:6379"),
		Username: web.AppConfig.DefaultString("redis::redis_username", ""),
		Password: web.AppConfig.DefaultString("redis::redis_password", ""),
	})
	err := cache.InitCache(cache.NewRedisCache(client))
	if err != nil {
		panic(err)
	}
}
