package Storage

import (
	"github.com/go-redis/redis"
	"golang.org/x/net/context"
)

func NewRedisClient(ctx context.Context) *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	conn.WithContext(ctx)
	return conn
}
