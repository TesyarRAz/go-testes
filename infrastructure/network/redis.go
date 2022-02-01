package network

import "github.com/go-redis/redis/v8"

type Redis struct {
	Client *redis.Client
}

func NewRedis(host, port, password string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	return &Redis{client}
}
