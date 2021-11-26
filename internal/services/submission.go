package services

import (
	"context"
	"log"
	"site/internal/datastruct"
	messagebroker "site/internal/message_broker"
	"site/internal/middleware"
	"site/internal/store"

	lru "github.com/hashicorp/golang-lru"
)

const (
	submissionsPerPage = 20
)

type SubmissionService interface {
	All(ctx context.Context, query *datastruct.SubmissionQuery) ([]*datastruct.Submission, error)
	ById(ctx context.Context, id int) (*datastruct.Submission, error)
	Create(ctx context.Context, submission *datastruct.Submission) error
	Update(ctx context.Context, submission *datastruct.Submission) error
	Delete(ctx context.Context, id int) error
}

type SubmissionServiceImpl struct {
	store  store.Store
	cache  *lru.TwoQueueCache
	broker messagebroker.MessageBroker
}

func NewSubmissionService(opts ...SubmissionServiceOption) SubmissionService {
	svc := &SubmissionServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (s SubmissionServiceImpl) All(ctx context.Context, query *datastruct.SubmissionQuery) ([]*datastruct.Submission, error) {
	if value, ok := s.cache.Get(query.Filter); ok {
		if submissions, ok := value.([]*datastruct.Submission); ok {
			log.Println("The result of cache", submissions)
			return submissions, nil
		}
	}
	query.Limit = submissionsPerPage
	query.Offset = (query.Page - 1) * submissionsPerPage
	submissions, err := s.store.Submissions().All(ctx, query)
	log.Println("Successfully cached!")
	s.cache.Add(query.Filter, submissions)
	return submissions, err
}

func (s SubmissionServiceImpl) ById(ctx context.Context, id int) (*datastruct.Submission, error) {
	value, ok := s.cache.Get(id)
	if ok {
		if submission, ok := value.(*datastruct.Submission); ok {
			log.Println("The result of cache", *submission)
			return submission, nil
		}
	}
	submission, err := s.store.Submissions().ById(ctx, id)
	log.Println("Successfully cached!")
	s.cache.Add(id, submission)
	return submission, err
}

func (s SubmissionServiceImpl) Create(ctx context.Context, submission *datastruct.Submission) error {
	// TODO: ASSIGN USER TO SUBMISSION
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}
	submission.AuthorHandle = user.Handle
	log.Println("Cache was purged")
	s.broker.Cache().Purge()
	return s.store.Submissions().Create(ctx, submission)
}

func (s SubmissionServiceImpl) Update(ctx context.Context, submission *datastruct.Submission) error {
	_, err := s.store.Submissions().ById(ctx, int(submission.Id))
	if err != nil {
		return err
	}
	log.Printf("Submission with %d was removed from cache\n", submission.Id)
	s.broker.Cache().Remove(submission.Id)
	return s.store.Submissions().Update(ctx, submission)
}

func (s SubmissionServiceImpl) Delete(ctx context.Context, id int) error {
	_, err := s.store.Submissions().ById(ctx, id)
	if err != nil {
		return err
	}
	s.broker.Cache().Remove(id)
	log.Printf("Submission with %d was removed from cache\n", id)
	return s.store.Submissions().Delete(ctx, id)
}
