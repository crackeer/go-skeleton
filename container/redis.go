package container

import (
	"time"

	redis "github.com/go-redis/redis/v8"
)

var (
	redisClients map[string]*redis.Client = make(map[string]*redis.Client)
)

func initRedis(list []RedisConfig) error {
	for _, conf := range list {
		client := redis.NewClient(&redis.Options{
			Addr:        conf.Host,
			Username:    conf.Username,
			Password:    conf.Password,
			DB:          conf.DB,
			DialTimeout: 10 * time.Second,
		})
		redisClients[conf.Name] = client
	}
	return nil
}
