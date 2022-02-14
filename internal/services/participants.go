package services

import (
	"context"
	"math/rand"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/store"
)

type ParticipantService interface {
	FindAll(ctx context.Context, req *dto.ParticipantFindAllRequest) ([]*datastruct.Participant, error)
	FindFriends(ctx context.Context, req *dto.ParticipantFindFriendsRequest) ([]*datastruct.Participant, error)
	Register(ctx context.Context, req *dto.ParticipantRegisterRequest) error
}

type ParticipantServiceImpl struct {
	store store.Store
}

var (
	rooms = 100
)

func NewParticipantRepository(opts ...ParticipantServiceOption) ParticipantService {
	svc := &ParticipantServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (p ParticipantServiceImpl) FindAll(ctx context.Context, req *dto.ParticipantFindAllRequest) ([]*datastruct.Participant, error) {
	req.Limit = usersPerPage
	req.Offset = (req.Page - 1) * usersPerPage

	return p.store.Participants().FindAll(ctx, req)
}

func (p ParticipantServiceImpl) FindFriends(ctx context.Context, req *dto.ParticipantFindFriendsRequest) ([]*datastruct.Participant, error) {
	req.Limit = usersPerPage
	req.Offset = (req.Page - 1) * usersPerPage

	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return nil, middleware.ErrNotAuthenticated
	}
	req.UserId = int(user.Id)

	return p.store.Participants().FindFriends(ctx, req)
}

func (p ParticipantServiceImpl) Register(ctx context.Context, req *dto.ParticipantRegisterRequest) error {
	if req.ParticipantType == "" {
		req.ParticipantType = dto.CONTESTANT.String()
	}
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}
	participant := &datastruct.Participant{
		UserId:          user.Id,
		ContestId:       int32(req.ContestId),
		ParticipantType: req.ParticipantType,
		Room:            int32(rand.Int()%rooms + 1),
	}
	return p.store.Participants().Create(ctx, participant)
}
