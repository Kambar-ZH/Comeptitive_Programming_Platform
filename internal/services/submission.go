package services

import (
	"context"
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
}

func NewSubmissionService(opts ...SubmissionServiceOption) SubmissionService {
	svc := &SubmissionServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (s SubmissionServiceImpl) All(ctx context.Context, query *datastruct.SubmissionQuery) ([]*datastruct.Submission, error) {
	query.Limit = submissionsPerPage
	query.Offset = (query.Page - 1) * submissionsPerPage
	return s.store.Submissions().All(ctx, query)
}

func (s SubmissionServiceImpl) ById(ctx context.Context, id int) (*datastruct.Submission, error) {
	return s.store.Submissions().ById(ctx, id)
}

func (s SubmissionServiceImpl) Create(ctx context.Context, submission *datastruct.Submission) error {
	// TODO: ASSIGN USER TO SUBMISSION
	user := middleware.UserFromCtx(ctx)
	submission.AuthorHandle = user.Handle
	return s.store.Submissions().Create(ctx, submission)
}

func (s SubmissionServiceImpl) Update(ctx context.Context, submission *datastruct.Submission) error {
	_, err := s.store.Submissions().ById(ctx, int(submission.Id))
	if err != nil {
		return err
	}
	return s.store.Submissions().Update(ctx, submission)
}

func (s SubmissionServiceImpl) Delete(ctx context.Context, id int) error {
	_, err := s.store.Submissions().ById(ctx, id)
	if err != nil {
		return err
	}
	return s.store.Submissions().Delete(ctx, id)
}