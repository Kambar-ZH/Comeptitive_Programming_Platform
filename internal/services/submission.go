package services

import (
	"context"
	"errors"
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

var (
	errSubmissionNotFound = errors.New("submission not found")
)

type SubmissionService interface {
	All(ctx context.Context, req *datastruct.SubmissionAllRequest) ([]*datastruct.Submission, error)
	ById(ctx context.Context, req *datastruct.SubmissionByIdRequest) (*datastruct.Submission, error)
	Create(ctx context.Context, submission *datastruct.Submission) error
	Update(ctx context.Context, req *datastruct.SubmissionUpdateRequest) error
	Delete(ctx context.Context, req *datastruct.SubmissionDeleteRequest) error
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

func (s SubmissionServiceImpl) All(ctx context.Context, req *datastruct.SubmissionAllRequest) ([]*datastruct.Submission, error) {
	if req.Filter != "" {
		if value, ok := s.cache.Get(req.Filter); ok {
			if submissions, ok := value.([]*datastruct.Submission); ok {
				log.Println("The result of cache", submissions)
				return submissions, nil
			}
		}
	}
	req.Limit = submissionsPerPage
	req.Offset = (req.Page - 1) * submissionsPerPage
	submissions, err := s.store.Submissions().All(ctx, req)
	if req.Filter != "" {
		log.Println("Successfully cached!")
		s.cache.Add(req.Filter, submissions)
	}
	return submissions, err
}

func (s SubmissionServiceImpl) ById(ctx context.Context, req *datastruct.SubmissionByIdRequest) (*datastruct.Submission, error) {
	value, ok := s.cache.Get(req.SubmissionId)
	if ok {
		if submission, ok := value.(*datastruct.Submission); ok {
			if submission.ContestId != req.ContestId {
				return nil, errSubmissionNotFound
			}
			return submission, nil
		}
	}
	submission, err := s.store.Submissions().ById(ctx, int(req.SubmissionId))
	if err != nil {
		return nil, err
	}
	if submission.ContestId != req.ContestId {
		return nil, errSubmissionNotFound
	}
	log.Println("Successfully cached!")
	s.cache.Add(req.SubmissionId, submission)
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

func (s SubmissionServiceImpl) Update(ctx context.Context, req *datastruct.SubmissionUpdateRequest) error {
	submission, err := s.store.Submissions().ById(ctx, int(req.Submission.Id))
	if err != nil {
		return err
	}
	if req.ContestId != submission.ContestId {
		return errSubmissionNotFound
	}
	log.Printf("Submission with [id=%d] was removed from cache\n", req.Submission.Id)
	s.broker.Cache().Remove(req.Submission.Id)
	return s.store.Submissions().Update(ctx, req.Submission)
}

func (s SubmissionServiceImpl) Delete(ctx context.Context, req *datastruct.SubmissionDeleteRequest) error {
	submission, err := s.store.Submissions().ById(ctx, int(req.SubmissionId))
	if err != nil {
		return err
	}
	if req.ContestId != submission.ContestId {
		return errSubmissionNotFound
	}
	s.broker.Cache().Remove(req.SubmissionId)
	log.Printf("Submission with [id=%d] was removed from cache\n", req.SubmissionId)
	return s.store.Submissions().Delete(ctx, int(req.SubmissionId))
}
