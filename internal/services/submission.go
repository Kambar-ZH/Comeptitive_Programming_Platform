package services

import (
	"context"
	"fmt"
	"log"
	"site/internal/cache"
	"site/internal/datastruct"
	"site/internal/middleware"
	"site/internal/store"
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
	store store.Store
	cache cache.Cache
}

func NewSubmissionService(opts ...SubmissionServiceOption) SubmissionService {
	svc := &SubmissionServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (s SubmissionServiceImpl) All(ctx context.Context, query *datastruct.SubmissionQuery) ([]*datastruct.Submission, error) {
	if query.Filter != "" {
		submissions, err := s.cache.Submissions().GetAll(query.Filter)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Redis ok!")
			return submissions, nil
		}
	}
	query.Limit = submissionsPerPage
	query.Offset = (query.Page - 1) * submissionsPerPage
	submissions, err := s.store.Submissions().All(ctx, query)
	if query.Filter != "" {
		log.Println("Saved to Redis!")
		if err := s.cache.Submissions().SetAll(query.Filter, submissions); err != nil {
			log.Println(err)
		}
	}
	return submissions, err
}

func (s SubmissionServiceImpl) ById(ctx context.Context, id int) (*datastruct.Submission, error) {
	submission, err := s.cache.Submissions().Get(fmt.Sprintf("%d", id))
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Redis ok!")
		return submission, nil
	}
	submission, err = s.store.Submissions().ById(ctx, id)
	if cacheErr := s.cache.Submissions().Set(fmt.Sprintf("%d", id), submission); cacheErr != nil {
		log.Println(cacheErr)
	}
	return submission, err
}

func (s SubmissionServiceImpl) Create(ctx context.Context, submission *datastruct.Submission) error {
	// TODO: ASSIGN USER TO SUBMISSION
	user := middleware.UserFromCtx(ctx)
	submission.AuthorHandle = user.Handle
	if cacheErr := s.cache.Submissions().Set(fmt.Sprintf("%d", submission.Id), submission); cacheErr != nil {
		log.Println(cacheErr)
	}
	return s.store.Submissions().Create(ctx, submission)
}

func (s SubmissionServiceImpl) Update(ctx context.Context, submission *datastruct.Submission) error {
	_, err := s.store.Submissions().ById(ctx, int(submission.Id))
	if err != nil {
		return err
	}
	if cacheErr := s.cache.Submissions().Set(fmt.Sprintf("%d", submission.Id), submission); cacheErr != nil {
		log.Println(cacheErr)
	}
	return s.store.Submissions().Update(ctx, submission)
}

func (s SubmissionServiceImpl) Delete(ctx context.Context, id int) error {
	_, err := s.store.Submissions().ById(ctx, id)
	if err != nil {
		return err
	}
	s.cache.Submissions().Del(fmt.Sprintf("%d", id))
	return s.store.Submissions().Delete(ctx, id)
}
