package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/go-redis/redis/v8"

	"test-task/shared/config"
	"test-task/shared/log"
)

type RedisService struct {
	client *redis.Client
	prefix string
}

var (
	pool  *redis.Client
	ponce sync.Once
)

func resetPool() {
	pool = nil
	ponce = sync.Once{}
}

// CreateConnection is made for establishing the redis connection
func CreateConnection(config config.IConfig) (*redis.Client, error) {
	var err error
	ponce.Do(func() {
		err = retry.Do(func() error {
			opts := redis.Options{
				Addr:     config.Redis().Host,
				DB:       config.Redis().Database,
				Password: config.Redis().Password,
				Username: config.Redis().Username,
			}

			pool = redis.NewClient(&opts)
			return nil
		},
			retry.Delay(1*time.Second))
		if err != nil {
			log.GetLog().Error("Redis connection error", err.Error())
		}
	})

	return pool, err
}

// GetConnection is made for getting the redis client for current connection
func GetConnection() (*redis.Client, error) {
	if pool == nil {
		return nil, fmt.Errorf("connection Pool has not been created")
	}
	return pool, nil
}

// ClosePool is used for closing the redis connection
func ClosePool() error {
	if pool == nil {
		return nil
	}
	return pool.Close()
}

// SetToken is a function used for setting the jwt token as redis key:value
func SetToken(ctx context.Context, id, token string, time time.Duration) error {
	redisConn, err := GetConnection()
	if err != nil {
		return err
	}
	err = redisConn.Set(ctx, id, token, time).Err()
	if err != nil {
		return fmt.Errorf("failed to set token in Redis: %w", err)
	}

	return nil
}

// GetToken is a function used for retrieving the jwt token value by redis key
func GetToken(ctx context.Context, id string) (string, error) {
	redisConn, err := GetConnection()
	if err != nil {
		return "", err
	}

	return redisConn.Get(ctx, id).Val(), nil
}
