package cache

import "site/internal/datastruct"

type SubmissionCache interface {
	Set(key string, value *datastruct.Submission) error
	Get(key string) (*datastruct.Submission, error)
	// Del(key string) error
	SetAll(key string, value []*datastruct.Submission)  error
	GetAll(key string) ([]*datastruct.Submission, error)
}