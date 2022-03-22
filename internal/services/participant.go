package services

import (
	"context"
	"math/rand"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/logger"
	"site/internal/middleware"
	"site/internal/store"
)

type ParticipantService interface {
	FindAll(ctx context.Context, req *dto.ParticipantFindAllRequest) ([]*datastruct.Participant, error)
	FindFriends(ctx context.Context, req *dto.ParticipantFindFriendsRequest) ([]*datastruct.Participant, error)
	Register(ctx context.Context, req *dto.ParticipantRegisterRequest) error
	Standings(ctx context.Context, req *dto.GetStandingsRequest) ([]*datastruct.StandingsRow, error)
}

type ParticipantServiceImpl struct {
	store store.Store
}

var (
	rooms = 100
)

func NewParticipantService(opts ...ParticipantServiceOption) ParticipantService {
	svc := &ParticipantServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (p ParticipantServiceImpl) Standings(ctx context.Context, req *dto.GetStandingsRequest) ([]*datastruct.StandingsRow, error) {
	req.Limit = usersPerPage
	req.Offset = (req.Page - 1) * usersPerPage
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return nil, middleware.ErrNotAuthenticated
	}
	req.UserId = user.Id
	standings, err := p.mapStandings(ctx, req)
	if err != nil {
		return nil, err
	}
	return standings, nil
}

func (p ParticipantServiceImpl) mapStandings(ctx context.Context, req *dto.GetStandingsRequest) ([]*datastruct.StandingsRow, error) {
	var standings = make([]*datastruct.StandingsRow, 0)
	participants, err := p.store.Participants().FindAll(ctx, &dto.ParticipantFindAllRequest{
		ContestId:  req.ContestId,
		Pagination: req.Pagination,
	})
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}
	problems, err := p.store.Problems().FindAll(ctx, &dto.ProblemFindAllRequest{
		ContestId:    req.ContestId,
		LanguageCode: middleware.LanguageCodeFromCtx(ctx),
	})
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	for _, participant := range participants {
		problemResults := make([]*datastruct.ProblemResult, 0)
		for _, problem := range problems {
			problemResult, err := p.store.ProblemResults().GetByProblemId(ctx, &dto.ProblemResultGetByProblemIdRequest{
				ProblemId: problem.Id,
			})
			if err != nil {
				logger.Logger.Error(err.Error())
				return nil, err
			}
			problemResults = append(problemResults, problemResult)
		}
		standings = append(standings, &datastruct.StandingsRow{
			Participant:    participant,
			ProblemResults: problemResults,
		})
	}
	return standings, nil
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
	req.UserId = user.Id

	return p.store.Participants().FindFriends(ctx, req)
}

func (p ParticipantServiceImpl) Register(ctx context.Context, req *dto.ParticipantRegisterRequest) error {
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}
	participant := &datastruct.Participant{
		UserId:          user.Id,
		ContestId:       int32(req.ContestId),
		ParticipantType: req.ParticipantType,
		Room:            int32(rand.Intn(rooms-1) + 1),
	}
	return p.store.Participants().Create(ctx, participant)
}
