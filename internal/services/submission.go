package services

import (
	"context"
	"fmt"
	"site/internal/grpc/api"
	"site/internal/store"
)

const (
	submissionsLimitPerPage = 20
)

type SubmissionService interface {
	All(ctx context.Context, req *api.AllSubmissionsRequest) (*api.SubmissionList, error)
	ById(ctx context.Context, req *api.SubmissionByIdRequest) (*api.Submission, error)
	ByAuthorHandle(ctx context.Context, req *api.SubmissionByHandleRequest) (*api.Submission, error)
	Create(ctx context.Context, submission *api.Submission) (*api.Submission, error)
	Update(ctx context.Context, submission *api.Submission) (*api.Submission, error) 
	Delete(ctx context.Context, req *api.DeleteSubmissionRequest) (*api.Empty, error)
}

type SubmissionServiceImpl struct {
	repo store.SubmissionRepository
}

func NewSubmissionService(opts ...SubmissionServiceOption) SubmissionService {
	svc := &SubmissionServiceImpl{}
	for _, v := range (opts) {
		v(svc)
	}
	return svc
}

func (s SubmissionServiceImpl) All(ctx context.Context, req *api.AllSubmissionsRequest) (*api.SubmissionList, error) {
	req.Limit = submissionsLimitPerPage
	return s.repo.All(ctx, req)
}

func (s SubmissionServiceImpl) ById(ctx context.Context, req *api.SubmissionByIdRequest) (*api.Submission, error) {
	return s.repo.ById(ctx, req)
}

func (s SubmissionServiceImpl) ByAuthorHandle(ctx context.Context, req *api.SubmissionByHandleRequest) (*api.Submission, error) {
	return s.repo.ByAuthorHandle(ctx, req)
}

func (s SubmissionServiceImpl) Create(ctx context.Context, submission *api.Submission) (*api.Submission, error) {
	return s.repo.Create(ctx, submission)
}

func (s SubmissionServiceImpl) Update(ctx context.Context, submission *api.Submission) (*api.Submission, error) {
	req := &api.SubmissionByIdRequest{
		Id: submission.Id,
	}
	_, err := s.repo.ById(ctx, req)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, submission)
}

func (s SubmissionServiceImpl) Delete(ctx context.Context, req *api.DeleteSubmissionRequest) (*api.Empty, error) {
	req2 := &api.SubmissionByIdRequest{
		Id: req.Id,
	}
	submission, err := s.repo.ById(ctx, req2)
	if err != nil {
		return nil, err
	}
	if submission.ContestId != req.ContestId {
		return nil, fmt.Errorf("contestId doesn't match to submissionId")
	}
	return s.repo.Delete(ctx, req)
}