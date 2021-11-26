package lru

import (
	"site/internal/cache"
	"time"
)

type Cache struct {
	submissions cache.SubmissionCache
}

func (cache *Cache) Connect(host string, db int, expires time.Duration) {}

func NewCache() cache.Cache {
	return &Cache{}
}
