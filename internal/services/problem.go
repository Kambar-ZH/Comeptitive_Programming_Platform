package services

import (
	"context"
	"errors"
	"site/internal/datastruct"
	"site/internal/store"
)

var (
	errProblemNotFound = errors.New("problem not found")
)

const (
	problemsPerPage = 20
)

type ProblemService interface {
	Problemset(ctx context.Context, query *datastruct.ProblemsetRequest) ([]*datastruct.Problem, error)
	All(ctx context.Context, req *datastruct.ProblemAllRequest) ([]*datastruct.Problem, error)
	ById(ctx context.Context, req *datastruct.ProblemByIdRequest) (*datastruct.Problem, error)
	Create(ctx context.Context, req *datastruct.ProblemCreateRequest) error
	Update(ctx context.Context, req *datastruct.ProblemUpdateRequest) error
	Delete(ctx context.Context, req *datastruct.ProblemDeleteRequest) error
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

func (p ProblemServiceImpl) Problemset(ctx context.Context, req *datastruct.ProblemsetRequest) ([]*datastruct.Problem, error) {
	req.Limit = problemsPerPage
	req.Offset = (req.Page - 1) * problemsPerPage
	return p.store.Problems().Problemset(ctx, req)
}

func (p ProblemServiceImpl) All(ctx context.Context, req *datastruct.ProblemAllRequest) ([]*datastruct.Problem, error) {
	req.Limit = problemsPerPage
	req.Offset = (req.Page - 1) * problemsPerPage
	return p.store.Problems().All(ctx, req)
}

func (p ProblemServiceImpl) ById(ctx context.Context, req *datastruct.ProblemByIdRequest) (*datastruct.Problem, error) {
	problem, err := p.store.Problems().ById(ctx, int(req.ProblemId))
	if err != nil {
		return nil, err
	}
	if problem.ContestId != req.ContestId {
		return nil, errProblemNotFound
	}
	return problem, nil
}

func (p ProblemServiceImpl) Create(ctx context.Context, req *datastruct.ProblemCreateRequest) error {
	req.Problem.ContestId = req.ContestId
	return p.store.Problems().Create(ctx, req.Problem)
}

func (p ProblemServiceImpl) Update(ctx context.Context, req *datastruct.ProblemUpdateRequest) error {
	// TODO: add admin permission
	problem, err := p.store.Problems().ById(ctx, int(req.Problem.Id))
	if err != nil {
		return err
	}
	if req.ContestId != problem.ContestId {
		return errProblemNotFound
	}
	return p.store.Problems().Update(ctx, req.Problem)
}

func (p ProblemServiceImpl) Delete(ctx context.Context, req *datastruct.ProblemDeleteRequest) error {
	// TODO: add admin permission
	problem, err := p.store.Problems().ById(ctx, int(req.ProblemId))
	if err != nil {
		return err
	}
	if req.ContestId != problem.ContestId {
		return errProblemNotFound
	}
	return p.store.Problems().Delete(ctx, int(req.ProblemId))
}
