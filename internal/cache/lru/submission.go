package lru

import (
	"fmt"
	"site/internal/cache"
	"site/internal/datastruct"

	lru "github.com/hashicorp/golang-lru"
)

type SubmissionCache struct {
	cache *lru.TwoQueueCache
}

func (cache *Cache) Submissions() cache.SubmissionCache {
	if cache.submissions == nil {
		cache.submissions = NewSubmissionCache()
	}
	return cache.submissions
}

func NewSubmissionCache() cache.SubmissionCache {
	cache, _ := lru.New2Q(100)
	return &SubmissionCache{
		cache: cache,
	}
}

func (s *SubmissionCache) Set(key interface{}, value *datastruct.Submission) error {
	s.cache.Add(key, value)
	return nil
}

func (s *SubmissionCache) Get(key interface{}) (*datastruct.Submission, error) {
	if value, ok := s.cache.Get(key); ok {
		if submission, ok := value.(*datastruct.Submission); ok {
			return submission, nil
		}
		return nil, fmt.Errorf("cannot parse the value with key - %v to *datastruct.Submission", key)
	}
	return nil, fmt.Errorf("no values found with given key")
}

func (s *SubmissionCache) SetAll(key interface{}, value []*datastruct.Submission) error {
	s.cache.Add(key, value)
	return nil
}

func (s *SubmissionCache) GetAll(key interface{}) ([]*datastruct.Submission, error) {
	if value, ok := s.cache.Get(key); ok {
		if submissions, ok := value.([]*datastruct.Submission); ok {
			return submissions, nil
		}
		return nil, fmt.Errorf("cannot parse the value with key - %v to []*datastruct.Submission", key)
	}
	return nil, fmt.Errorf("no values found with given key")
}

func (s *SubmissionCache) Del(key interface{}) error {
	s.cache.Remove(key)
	return nil
}