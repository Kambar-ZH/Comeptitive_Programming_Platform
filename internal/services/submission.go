package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"
)

const (
	submissionsPerPage = 20
)

type SubmissionService interface {
	All(ctx context.Context, page int) ([]*datastruct.Submission, error)
	ByAuthorHandle(ctx context.Context, handle string) ([]*datastruct.Submission, error)
	ByContestId(ctx context.Context, contestId int) ([]*datastruct.Submission, error)
	ById(ctx context.Context, id int) (*datastruct.Submission, error)
	Create(ctx context.Context, submission *datastruct.Submission) error
	Update(ctx context.Context, submission *datastruct.Submission) error
	Delete(ctx context.Context, id int) error
}

type SubmissionServiceImpl struct {
	repo store.SubmissionRepository
}

func NewSubmissionService(opts ...SubmissionServiceOption) SubmissionService {
	svc := &SubmissionServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (s SubmissionServiceImpl) ByContestId(ctx context.Context, contestId int) ([]*datastruct.Submission, error) {
	return s.repo.ByContestId(ctx, contestId)
}

func (s SubmissionServiceImpl) All(ctx context.Context, page int) ([]*datastruct.Submission, error) {
	return s.repo.All(ctx, (page-1)*submissionsPerPage, submissionsPerPage)
}

func (s SubmissionServiceImpl) ById(ctx context.Context, id int) (*datastruct.Submission, error) {
	return s.repo.ById(ctx, id)
}

func (s SubmissionServiceImpl) ByAuthorHandle(ctx context.Context, handle string) ([]*datastruct.Submission, error) {
	return s.repo.ByAuthorHandle(ctx, handle)
}

func (s SubmissionServiceImpl) Create(ctx context.Context, submission *datastruct.Submission) error {
	return s.repo.Create(ctx, submission)
}

func (s SubmissionServiceImpl) Update(ctx context.Context, submission *datastruct.Submission) error {
	_, err := s.repo.ById(ctx, int(submission.Id))
	if err != nil {
		return err
	}
	return s.repo.Update(ctx, submission)
}

func (s SubmissionServiceImpl) Delete(ctx context.Context, id int) error {
	_, err := s.repo.ById(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
