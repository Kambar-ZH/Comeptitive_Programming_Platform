package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/logger"
	message_broker "site/internal/message_broker"
	"site/internal/middleware"
	"site/internal/store"
)

const (
	contestPerPage = 10
)

type ContestService interface {
	FindAll(ctx context.Context, query *dto.ContestFindAllRequest) ([]*datastruct.Contest, error)
	FindByTimeInterval(ctx context.Context, req *dto.ContestFindByTimeIntervalRequest) ([]*datastruct.Contest, error)
	GetById(ctx context.Context, req *dto.ContestGetByIdRequest) (*datastruct.Contest, error)
	Create(ctx context.Context, contest *datastruct.Contest) error
	Update(ctx context.Context, contest *datastruct.Contest) error
	Delete(ctx context.Context, id int) error
}

type ContestServiceImpl struct {
	store  store.Store
	broker message_broker.MessageBroker
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

func (c ContestServiceImpl) FindByTimeInterval(ctx context.Context, req *dto.ContestFindByTimeIntervalRequest) ([]*datastruct.Contest, error) {
	return c.store.Contests().FindByTimeInterval(ctx, req)
}

func (c ContestServiceImpl) GetById(ctx context.Context, req *dto.ContestGetByIdRequest) (*datastruct.Contest, error) {
	return c.store.Contests().GetById(ctx, req)
}

func (c ContestServiceImpl) Create(ctx context.Context, contest *datastruct.Contest) error {
	err := c.store.Contests().Create(ctx, contest)
	if err != nil {
		return err
	}
	if err := c.broker.Contest().CreateContest(contest); err != nil {
		logger.Logger.Error(err.Error())
	}
	return nil
}

func (c ContestServiceImpl) Update(ctx context.Context, contest *datastruct.Contest) error {
	_, err := c.store.Contests().GetById(ctx, &dto.ContestGetByIdRequest{
		ContestId:    contest.Id,
		LanguageCode: middleware.LanguageCodeFromCtx(ctx),
	})
	if err != nil {
		return err
	}
	return c.store.Contests().Update(ctx, contest)
}

func (c ContestServiceImpl) Delete(ctx context.Context, id int) error {
	_, err := c.store.Contests().GetById(ctx, &dto.ContestGetByIdRequest{
		ContestId:    int32(id),
		LanguageCode: middleware.LanguageCodeFromCtx(ctx),
	})
	if err != nil {
		return err
	}
	return c.store.Contests().Delete(ctx, id)
}
