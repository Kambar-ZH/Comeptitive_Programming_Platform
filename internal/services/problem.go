package services

import (
	"context"
	"errors"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/store"
)

var (
	errProblemNotFound = errors.New("problem not found")
)

const (
	problemsPerPage = 20
)

type ProblemService interface {
	Problemset(ctx context.Context, req *dto.ProblemsetRequest) ([]*datastruct.Problem, error)
	FindAll(ctx context.Context, req *dto.ProblemFindAllRequest) ([]*datastruct.Problem, error)
	GetById(ctx context.Context, req *dto.ProblemGetByIdRequest) (*datastruct.Problem, error)
	Create(ctx context.Context, req *dto.ProblemCreateRequest) error
	Update(ctx context.Context, req *dto.ProblemUpdateRequest) error
	Delete(ctx context.Context, req *dto.ProblemDeleteRequest) error
}

type ProblemServiceImpl struct {
	store store.Store
}

func NewProblemService(opts ...ProblemServiceOption) ProblemService {
	svc := &ProblemServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (p ProblemServiceImpl) Problemset(ctx context.Context, req *dto.ProblemsetRequest) ([]*datastruct.Problem, error) {
	req.Limit = problemsPerPage
	req.Offset = (req.Page - 1) * problemsPerPage
	return p.store.Problems().Problemset(ctx, req)
}

func (p ProblemServiceImpl) FindAll(ctx context.Context, req *dto.ProblemFindAllRequest) ([]*datastruct.Problem, error) {
	req.Limit = problemsPerPage
	req.Offset = (req.Page - 1) * problemsPerPage
	return p.store.Problems().FindAll(ctx, req)
}

func (p ProblemServiceImpl) GetById(ctx context.Context, req *dto.ProblemGetByIdRequest) (*datastruct.Problem, error) {
	problem, err := p.store.Problems().GetById(ctx, int(req.ProblemId))
	if err != nil {
		return nil, err
	}
	if problem.ContestId != req.ContestId {
		return nil, errProblemNotFound
	}
	return problem, nil
}

func (p ProblemServiceImpl) Create(ctx context.Context, req *dto.ProblemCreateRequest) error {
	req.Problem.ContestId = req.ContestId
	return p.store.Problems().Create(ctx, req.Problem)
}

func (p ProblemServiceImpl) Update(ctx context.Context, req *dto.ProblemUpdateRequest) error {
	// TODO: add admin permission
	problem, err := p.store.Problems().GetById(ctx, int(req.Problem.Id))
	if err != nil {
		return err
	}
	if req.ContestId != problem.ContestId {
		return errProblemNotFound
	}
	return p.store.Problems().Update(ctx, req.Problem)
}

func (p ProblemServiceImpl) Delete(ctx context.Context, req *dto.ProblemDeleteRequest) error {
	// TODO: add admin permission
	problem, err := p.store.Problems().GetById(ctx, int(req.ProblemId))
	if err != nil {
		return err
	}
	if req.ContestId != problem.ContestId {
		return errProblemNotFound
	}
	return p.store.Problems().Delete(ctx, int(req.ProblemId))
}
