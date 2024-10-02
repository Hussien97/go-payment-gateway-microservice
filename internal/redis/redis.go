package redis

import (
	"context"
	"log"
	"os"

	"payment-gateway/internal/resilience"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

// initializes the Redis client
func InitRedis() {
	var err error

	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test the connection with resilience's service retry logic
	err = resilience.RetryOperation(func() error {
		return rdb.Ping(ctx).Err()
	}, 5)

	if err != nil {
		log.Fatalf("Could not connect to Redis after retries: %v", err)
	}

	log.Println("Connected to Redis successfully.")
}

// sets the transaction status in Redis made it with no expiration but can be improved later by expiring it once tranasction status become completed
func SetTransactionStatus(transactionID string, status string) {
	err := rdb.Set(ctx, transactionID, status, 0).Err()
	if err != nil {
		log.Printf("Could not set transaction status in Redis: %v", err)
	}
}

// retrieves the transaction status from Redis
func GetTransactionStatus(transactionID string) (string, error) {
	status, err := rdb.Get(ctx, transactionID).Result()
	if err != nil {
		return "", err
	}
	return status, nil
}
