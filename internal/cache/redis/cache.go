package redis

import (
	"site/internal/cache"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client *redis.Client

	submissions cache.SubmissionCache
}

func (cache *Cache) Connect(host string, db int, expires time.Duration) {
	cache.client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       db,
	})
}

func NewCache() cache.Cache {
	return &Cache{}
}
