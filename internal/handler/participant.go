package handler

import (
	"net/http"
	"site/internal/dto"
	"site/internal/logger"
	"site/internal/middleware"
	"site/internal/services"
	"strconv"

	"github.com/go-chi/chi"
)

type ParticipantHandler struct {
	service services.ParticipantService
}

func NewParticipantHandler(opts ...ParticipantHandlerOption) *ParticipantHandler {
	ph := &ParticipantHandler{}
	for _, v := range opts {
		v(ph)
	}
	return ph
}

func (ph *ParticipantHandler) Standings(w http.ResponseWriter, r *http.Request) {
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	standings, err := ph.service.Standings(r.Context(), &dto.GetStandingsRequest{
		ContestId: int32(contestId),
		Pagination: dto.Pagination{
			Page:   r.Context().Value(middleware.CtxKeyPage).(int32),
			Filter: r.Context().Value(middleware.CtxKeyFilter).(string),
		},
	})

	Respond(w, r, http.StatusOK, standings)

}

func (ph *ParticipantHandler) Register(w http.ResponseWriter, r *http.Request) {
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	participantType := r.URL.Query().Get("participant_type")
	req := &dto.ParticipantRegisterRequest{
		ContestId:       int32(contestId),
		ParticipantType: participantType,
	}
	err = ph.service.Register(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, nil)
}

func (ph *ParticipantHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		logger.Logger.Error(err.Error())
		return
	}

	req := &dto.ParticipantFindAllRequest{
		ContestId: int32(contestId),
		Pagination: dto.Pagination{
			Page:   r.Context().Value(middleware.CtxKeyPage).(int32),
			Filter: r.Context().Value(middleware.CtxKeyFilter).(string),
		},
	}

	users, err := ph.service.FindAll(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, users)
}

func (ph *ParticipantHandler) FindFriends(w http.ResponseWriter, r *http.Request) {
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		logger.Logger.Error(err.Error())
		return
	}

	req := &dto.ParticipantFindFriendsRequest{
		ContestId: int32(contestId),
		Pagination: dto.Pagination{
			Page:   r.Context().Value(middleware.CtxKeyPage).(int32),
			Filter: r.Context().Value(middleware.CtxKeyFilter).(string),
		},
	}

	users, err := ph.service.FindFriends(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, users)
}
