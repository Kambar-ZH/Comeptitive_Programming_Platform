package redis

import (
	"context"
	"encoding/json"
	"site/internal/datastruct"
	"site/internal/cache"
	"time"

	"github.com/go-redis/redis/v8"
)

type SubmissionCache struct {
	client *redis.Client
	expires time.Duration
}


func (cache *Cache) Submissions() cache.SubmissionCache {
	if cache.submissions == nil {
		cache.submissions = NewSubmissionCache(cache.client)
	}
	return cache.submissions
}

func NewSubmissionCache(client *redis.Client) cache.SubmissionCache {
	return &SubmissionCache{
		client: client,
		expires: 10,
	}
}

func (cache *SubmissionCache) Set(key string, value *datastruct.Submission) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cache.client.Set(context.Background(), key, json, cache.expires*time.Second)
	return nil
}

func (cache *SubmissionCache) Get(key string) (*datastruct.Submission, error) {
	val, err := cache.client.Get(context.Background(), key).Result()
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

func (cache *SubmissionCache) SetAll(key string, value []*datastruct.Submission) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cache.client.Set(context.Background(), key+"all_users", json, cache.expires*time.Second)
	return nil
}

func (cache *SubmissionCache) GetAll(key string) ([]*datastruct.Submission, error) {
	val, err := cache.client.Get(context.Background(), key+"all_users").Result()
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

func (cache *SubmissionCache) Del(key string) error {
	return cache.client.Del(context.Background(), key).Err()
}