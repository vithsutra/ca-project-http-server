package redisqueue

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueueRepo struct {
	client *redis.Client
}

func NewRedisQueueRepo(client *redis.Client) *RedisQueueRepo {
	return &RedisQueueRepo{
		client: client,
	}
}

func (r *RedisQueueRepo) SendEmail(data []byte) error {
	const queueName = "email_queue"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("📤 Sending data to Redis list:", string(data))

	// ✅ Check Redis connection health
	if err := r.client.Ping(ctx).Err(); err != nil {
		log.Println("❌ Redis ping failed before LPUSH:", err)
		return err
	}

	// ✅ Push data to Redis list (queue)
	if err := r.client.LPush(ctx, queueName, string(data)).Err(); err != nil {
		log.Println("❌ Redis LPUSH failed:", err)
		return err
	}

	log.Println("✅ Successfully pushed email job to Redis list")
	return nil
}
