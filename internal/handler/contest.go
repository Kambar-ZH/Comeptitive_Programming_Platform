package handler

import (
	"encoding/json"
	"net/http"
	"site/internal/consts"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/services"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type ContestHander struct {
	service services.ContestService
}

func NewContestHandler(opts ...ContestHandlerOption) *ContestHander {
	ch := &ContestHander{}
	for _, v := range opts {
		v(ch)
	}
	return ch
}

func (ch *ContestHander) Create(w http.ResponseWriter, r *http.Request) {
	contest := new(datastruct.Contest)
	if err := json.NewDecoder(r.Body).Decode(contest); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	err := ch.service.Create(r.Context(), contest)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (ch *ContestHander) FindAll(w http.ResponseWriter, r *http.Request) {
	req := &dto.ContestFindAllRequest{
		Pagination: dto.Pagination{
			Page:   r.Context().Value(middleware.CtxKeyPage).(int32),
			Filter: r.Context().Value(middleware.CtxKeyFilter).(string),
		},
		LanguageCode: middleware.LanguageCodeFromCtx(r.Context()),
	}

	contests, err := ch.service.FindAll(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, contests)
}

func (ch *ContestHander) FindByTimeInterval(w http.ResponseWriter, r *http.Request) {
	timeFromStr := r.URL.Query().Get("time_from")
	timeToStr := r.URL.Query().Get("time_to")

	timeFrom, err := time.Parse(time.RFC3339, timeFromStr)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	timeTo, err := time.Parse(time.RFC3339, timeToStr)
	if err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	contests, err := ch.service.FindByTimeInterval(r.Context(), &dto.ContestFindByTimeIntervalRequest{
		StartTime:    timeFrom,
		EndTime:      timeTo,
		LanguageCode: consts.EN,
	})

	Respond(w, r, http.StatusOK, contests)
}

func (ch *ContestHander) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "contestId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	contest, err := ch.service.GetById(r.Context(), &dto.ContestGetByIdRequest{
		ContestId:    int32(id),
		LanguageCode: middleware.LanguageCodeFromCtx(r.Context()),
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, contest)
}

func (ch *ContestHander) Update(w http.ResponseWriter, r *http.Request) {
	contest := new(datastruct.Contest)
	if err := json.NewDecoder(r.Body).Decode(contest); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	err := ch.service.Update(r.Context(), contest)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusAccepted, nil)
}

func (ch *ContestHander) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := ch.service.Delete(r.Context(), id); err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, nil)
}
