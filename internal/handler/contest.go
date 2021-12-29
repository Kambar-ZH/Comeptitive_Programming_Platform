package handler

import (
	"encoding/json"
	"net/http"
	"site/internal/datastruct"
	"site/internal/middleware"
	"site/internal/services"
	"strconv"

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

func (ch *ContestHander) All(w http.ResponseWriter, r *http.Request) {
	req := &datastruct.ContestAllRequest{
		Page:   r.Context().Value(middleware.CtxKeyPage).(int32),
		Filter: r.Context().Value(middleware.CtxKeyFilter).(string),
	}

	contests, err := ch.service.All(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, contests)
}

func (ch *ContestHander) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "contestId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	contest, err := ch.service.ById(r.Context(), id)
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