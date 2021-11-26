package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"site/internal/cache"
	"site/internal/datastruct"
	"time"

	"github.com/go-redis/redis/v8"
)

type SubmissionCache struct {
	client  *redis.Client
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
		client:  client,
		expires: 10,
	}
}

func (s *SubmissionCache) Set(key interface{}, value *datastruct.Submission) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	s.client.Set(context.Background(), fmt.Sprintf("%v", key), json, s.expires*time.Second)
	return nil
}

func (s *SubmissionCache) Get(key interface{}) (*datastruct.Submission, error) {
	val, err := s.client.Get(context.Background(), fmt.Sprintf("%v", key)).Result()
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

func (s *SubmissionCache) SetAll(key interface{}, value []*datastruct.Submission) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	s.client.Set(context.Background(), fmt.Sprintf("all_users%v", key), json, s.expires*time.Second)
	return nil
}

func (s *SubmissionCache) GetAll(key interface{}) ([]*datastruct.Submission, error) {
	val, err := s.client.Get(context.Background(), fmt.Sprintf("all_users%v", key)).Result()
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

func (s *SubmissionCache) Del(key interface{}) error {
	return s.client.Del(context.Background(), fmt.Sprintf("%v", key)).Err()
}
