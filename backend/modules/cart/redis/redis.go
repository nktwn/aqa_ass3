package redis

import (
	"github.com/go-redis/redis/v8"
)

type cartRedis struct {
	redisCl *redis.Client
}

func NewCartRedis(redisCl *redis.Client) *cartRedis {
	return &cartRedis{
		redisCl: redisCl,
	}
}
