package cache

import (
	"site/internal/datastruct"
	"time"
)

type Cache interface {
	Submissions() SubmissionCache
	Connect(host string, db int, expires time.Duration)
}

type SubmissionCache interface {
	Set(key interface{}, value *datastruct.Submission) error
	Get(key interface{}) (*datastruct.Submission, error)
	Del(key interface{}) error
	SetAll(key interface{}, value []*datastruct.Submission) error
	GetAll(key interface{}) ([]*datastruct.Submission, error)
}
