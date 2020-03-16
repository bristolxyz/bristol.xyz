package clients

import (
	"github.com/go-redis/redis/v7"
	"os"
)

// RedisClient is the Redis client we are using.
var RedisClient *redis.Client

// CreateRedisClient is used to create the Redis client.
func CreateRedisClient() error {
	Host := os.Getenv("REDIS_HOST")
	if Host == "" {
		Host = "localhost:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Host,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	_, err := RedisClient.Ping().Result()
	return err
}
