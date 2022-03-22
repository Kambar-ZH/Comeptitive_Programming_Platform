package services

import (
	"context"
	"errors"
	"site/internal/consts"
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
	problemset, err := p.store.Problems().Problemset(ctx, req)
	if err != nil {
		return nil, err
	}
	for i := range problemset {
		tags, err := p.store.Tags().GetByProblemId(ctx, int(problemset[i].Id))
		if err != nil {
			return nil, err
		}
		problemset[i].Tags = tags
	}
	return problemset, nil
}

func (p ProblemServiceImpl) FindAll(ctx context.Context, req *dto.ProblemFindAllRequest) ([]*datastruct.Problem, error) {
	problems, err := p.store.Problems().FindAll(ctx, req)
	if err != nil {
		return nil, err
	}
	for i := range problems {
		tags, err := p.store.Tags().GetByProblemId(ctx, int(problems[i].Id))
		if err != nil {
			return nil, err
		}
		problems[i].Tags = tags
	}
	return problems, nil
}

func (p ProblemServiceImpl) GetById(ctx context.Context, req *dto.ProblemGetByIdRequest) (*datastruct.Problem, error) {
	problem, err := p.store.Problems().GetById(ctx, req)
	if err != nil {
		return nil, err
	}
	if problem.ContestId != req.ContestId {
		return nil, errProblemNotFound
	}
	tags, err := p.store.Tags().GetByProblemId(ctx, int(problem.Id))
	if err != nil {
		return nil, err
	}
	problem.Tags = tags
	return problem, nil
}

func (p ProblemServiceImpl) Create(ctx context.Context, req *dto.ProblemCreateRequest) error {
	req.Problem.ContestId = req.ContestId
	return p.store.Problems().Create(ctx, req.Problem)
}

func (p ProblemServiceImpl) Update(ctx context.Context, req *dto.ProblemUpdateRequest) error {
	// TODO: add admin permission
	problem, err := p.store.Problems().GetById(ctx, &dto.ProblemGetByIdRequest{
		ProblemId: req.Problem.Id,
		ContestId: req.ContestId,
	})
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
	problem, err := p.store.Problems().GetById(ctx, &dto.ProblemGetByIdRequest{
		ProblemId:    req.ProblemId,
		ContestId:    req.ContestId,
		LanguageCode: consts.EN,
	})
	if err != nil {
		return err
	}
	if req.ContestId != problem.ContestId {
		return errProblemNotFound
	}
	return p.store.Problems().Delete(ctx, int(req.ProblemId))
}
