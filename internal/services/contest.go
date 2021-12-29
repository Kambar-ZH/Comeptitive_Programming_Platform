package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"
)

const (
	contestPerPage = 10
)

type ContestService interface {
	All(ctx context.Context, query *datastruct.ContestAllRequest) ([]*datastruct.Contest, error)
	ById(ctx context.Context, id int) (*datastruct.Contest, error)
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

func (c ContestServiceImpl) All(ctx context.Context, req *datastruct.ContestAllRequest) ([]*datastruct.Contest, error) {
	req.Limit = contestPerPage
	req.Offset = (req.Page - 1) * contestPerPage
	return c.store.Contests().All(ctx, req)
}

func (c ContestServiceImpl) ById(ctx context.Context, id int) (*datastruct.Contest, error) {
	return c.store.Contests().ById(ctx, id)
}

func (c ContestServiceImpl) Create(ctx context.Context, contest *datastruct.Contest) error {
	return c.store.Contests().Create(ctx, contest)
}

func (c ContestServiceImpl) Update(ctx context.Context, contest *datastruct.Contest) error {
	_, err := c.store.Contests().ById(ctx, int(contest.Id))
	if err != nil {
		return err
	}
	return c.store.Contests().Update(ctx, contest)
}

func (c ContestServiceImpl) Delete(ctx context.Context, id int) error {
	_, err := c.store.Contests().ById(ctx, id)
	if err != nil {
		return err
	}
	return c.store.Contests().Delete(ctx, id)
}