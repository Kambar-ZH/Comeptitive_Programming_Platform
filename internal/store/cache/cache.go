package cache

import (
	"context"
	"encoding/json"
	"site/internal/datastruct"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, expires time.Duration) SubmissionCache {
	return &RedisCache{
		host:    host,
		db:      db,
		expires: expires,
	}
}

func (cache *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) Set(key string, value *datastruct.Submission) error {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	client.Set(context.Background(), key, json, cache.expires*time.Second)
	return nil
}

func (cache *RedisCache) Get(key string) (*datastruct.Submission, error) {
	client := cache.getClient()

	val, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	submission := &datastruct.Submission{}
	err = json.Unmarshal([]byte(val), submission)
	if err != nil {
		return nil, err
	}
	return submission, nil
}

func (cache *RedisCache) SetAll(key string, value []*datastruct.Submission) error {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	client.Set(context.Background(), key + "all_users", json, cache.expires*time.Second)
	return nil
}

func (cache *RedisCache) GetAll(key string) ([]*datastruct.Submission, error) {
	client := cache.getClient()

	val, err := client.Get(context.Background(), key + "all_users").Result()
	if err != nil {
		return nil, err
	}

	submissions := make([]*datastruct.Submission, 0)
	err = json.Unmarshal([]byte(val), &submissions)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}