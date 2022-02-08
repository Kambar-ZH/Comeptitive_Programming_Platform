package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/store"
)

const (
	contestPerPage = 10
)

type ContestService interface {
	FindAll(ctx context.Context, query *dto.ContestFindAllRequest) ([]*datastruct.Contest, error)
	GetById(ctx context.Context, id int) (*datastruct.Contest, error)
	Create(ctx context.Context, contest *datastruct.Contest) error
	Update(ctx context.Context, contest *datastruct.Contest) error
	Delete(ctx context.Context, id int) error
}

type ContestServiceImpl struct {
	store store.Store
}

func NewContestService(opts ...ContestServiceOption) ContestService {
	svc := &ContestServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (c ContestServiceImpl) FindAll(ctx context.Context, req *dto.ContestFindAllRequest) ([]*datastruct.Contest, error) {
	req.Limit = contestPerPage
	req.Offset = (req.Page - 1) * contestPerPage
	return c.store.Contests().FindAll(ctx, req)
}

func (c ContestServiceImpl) GetById(ctx context.Context, id int) (*datastruct.Contest, error) {
	return c.store.Contests().GetById(ctx, id)
}

func (c ContestServiceImpl) Create(ctx context.Context, contest *datastruct.Contest) error {
	return c.store.Contests().Create(ctx, contest)
}

func (c ContestServiceImpl) Update(ctx context.Context, contest *datastruct.Contest) error {
	_, err := c.store.Contests().GetById(ctx, int(contest.Id))
	if err != nil {
		return err
	}
	return c.store.Contests().Update(ctx, contest)
}

func (c ContestServiceImpl) Delete(ctx context.Context, id int) error {
	_, err := c.store.Contests().GetById(ctx, id)
	if err != nil {
		return err
	}
	return c.store.Contests().Delete(ctx, id)
}