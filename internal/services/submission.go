package services

import (
	"context"
	"errors"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/logger"
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
	All(ctx context.Context, req *dto.SubmissionFindAllRequest) ([]*datastruct.Submission, error)
	ById(ctx context.Context, req *dto.SubmissionGetByIdRequest) (*datastruct.Submission, error)
	Create(ctx context.Context, req *dto.SubmissionCreateRequest) error
	Update(ctx context.Context, req *dto.SubmissionUpdateRequest) error
	Delete(ctx context.Context, req *dto.SubmissionDeleteRequest) error
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

func (s SubmissionServiceImpl) All(ctx context.Context, req *dto.SubmissionFindAllRequest) ([]*datastruct.Submission, error) {
	if req.Filter != "" {
		if value, ok := s.cache.Get(req.Filter); ok {
			if submissions, ok := value.([]*datastruct.Submission); ok {
				logger.Logger.Sugar().Debugf("cache returned", submissions)
				return submissions, nil
			}
		}
	}
	req.Limit = submissionsPerPage
	req.Offset = (req.Page - 1) * submissionsPerPage
	submissions, err := s.store.Submissions().FindAll(ctx, req)
	if req.Filter != "" {
		logger.Logger.Sugar().Debugf("query cached", submissions)
		s.cache.Add(req.Filter, submissions)
	}
	return submissions, err
}

func (s SubmissionServiceImpl) ById(ctx context.Context, req *dto.SubmissionGetByIdRequest) (*datastruct.Submission, error) {
	value, ok := s.cache.Get(req.SubmissionId)
	if ok {
		if submission, ok := value.(*datastruct.Submission); ok {
			if submission.ContestId != req.ContestId {
				return nil, errSubmissionNotFound
			}
			return submission, nil
		}
	}
	submission, err := s.store.Submissions().GetById(ctx, int(req.SubmissionId))
	if err != nil {
		return nil, err
	}
	if submission.ContestId != req.ContestId {
		return nil, errSubmissionNotFound
	}
	logger.Logger.Sugar().Debugf("query cached", submission)
	s.cache.Add(req.SubmissionId, submission)
	return submission, err
}

func (s SubmissionServiceImpl) Create(ctx context.Context, req *dto.SubmissionCreateRequest) error {
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}
	req.Submission.UserId = user.Id
	req.Submission.ContestId = req.ContestId
	logger.Logger.Debug("cache was purged")
	s.broker.Cache().Purge()
	return s.store.Submissions().Create(ctx, req.Submission)
}

func (s SubmissionServiceImpl) Update(ctx context.Context, req *dto.SubmissionUpdateRequest) error {
	submission, err := s.store.Submissions().GetById(ctx, int(req.Submission.Id))
	if err != nil {
		return err
	}
	if req.ContestId != submission.ContestId {
		return errSubmissionNotFound
	}
	logger.Logger.Sugar().Debugf("submission with [id=%d] was removed from cache", req.Submission.Id)
	s.broker.Cache().Remove(req.Submission.Id)
	return s.store.Submissions().Update(ctx, req.Submission)
}

func (s SubmissionServiceImpl) Delete(ctx context.Context, req *dto.SubmissionDeleteRequest) error {
	submission, err := s.store.Submissions().GetById(ctx, int(req.SubmissionId))
	if err != nil {
		return err
	}
	if req.ContestId != submission.ContestId {
		return errSubmissionNotFound
	}
	s.broker.Cache().Remove(req.SubmissionId)
	logger.Logger.Sugar().Debugf("submission with [id=%d] was removed from cache", req.SubmissionId)
	return s.store.Submissions().Delete(ctx, int(req.SubmissionId))
}
