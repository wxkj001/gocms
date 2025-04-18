package app

import (
	"gocms/cache"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var cacheModule = fx.Module("cacheModule", fx.Provide(NewCache))

func NewCache(config *viper.Viper) cache.Cache {
	switch config.GetString("cache.driver") {
	case "redis":
		return cache.NewRedisCache(redis.NewClient(&redis.Options{
			Addr:       config.GetString("cache.host") + ":" + config.GetString("cache.port"),
			Password:   config.GetString("cache.password"), // 没有密码，默认值
			DB:         config.GetInt("cache.db"),          // 默认DB 0
			MaxRetries: config.GetInt("cache.maxIdle"),
		}))
	}
	return nil
}
