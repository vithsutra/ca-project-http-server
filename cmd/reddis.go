package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueueConnection struct {
	Client *redis.Client
}

func ConnectToRedisQueue() *RedisQueueConnection {
	url := "redis://service.email.vithsutra.com:6379"

	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(fmt.Errorf("failed to parse REDIS_URL: %v", err))
	}

	opt.DialTimeout = 3 * time.Second
	opt.ReadTimeout = 2 * time.Second
	opt.WriteTimeout = 2 * time.Second

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("failed to connect to Redis: %v", err))
	}

	return &RedisQueueConnection{Client: client}
}
