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

type ProblemHandler struct {
	service services.ProblemService
}

func NewProblemHandler(opts ...ProblemHandlerOption) *ProblemHandler {
	ph := &ProblemHandler{}
	for _, v := range opts {
		v(ph)
	}
	return ph
}

func (ph *ProblemHandler) Create(w http.ResponseWriter, r *http.Request) {
	problem := new(datastruct.Problem)
	if err := json.NewDecoder(r.Body).Decode(problem); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	err = ph.service.Create(r.Context(), &datastruct.ProblemCreateRequest{
		Problem:   problem,
		ContestId: int32(contestId),
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, nil)
}

func (ph *ProblemHandler) Problemset(w http.ResponseWriter, r *http.Request) {

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	minDifficultyStr := r.URL.Query().Get("minDifficulty")
	minDifficulty, err := strconv.Atoi(minDifficultyStr)
	if err != nil {
		minDifficulty = 0
	}

	maxDifficultyStr := r.URL.Query().Get("maxDifficulty")
	maxDifficulty, err := strconv.Atoi(maxDifficultyStr)
	if err != nil {
		maxDifficulty = 5000
	}

	filterTag := r.URL.Query().Get("filterTag")

	req := &datastruct.ProblemsetRequest{
		Page:          int32(page),
		MinDifficulty: int32(minDifficulty),
		MaxDifficulty: int32(maxDifficulty),
		FilterTag:     filterTag,
	}

	problems, err := ph.service.Problemset(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, problems)
}

func (ph *ProblemHandler) All(w http.ResponseWriter, r *http.Request) {

	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	req := &datastruct.ProblemAllRequest{
		Page:      r.Context().Value(middleware.CtxKeyPage).(int32),
		ContestId: int32(contestId),
	}

	problems, err := ph.service.All(r.Context(), req)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, problems)
}

func (ph *ProblemHandler) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	submission, err := ph.service.ById(r.Context(), &datastruct.ProblemByIdRequest{
		ProblemId: int32(id),
		ContestId: int32(contestId),
	})
	if err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, submission)
}

func (ph *ProblemHandler) Update(w http.ResponseWriter, r *http.Request) {
	problem := new(datastruct.Problem)
	if err := json.NewDecoder(r.Body).Decode(problem); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	err = ph.service.Update(r.Context(), &datastruct.ProblemUpdateRequest{
		Problem:   problem,
		ContestId: int32(contestId),
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusAccepted, nil)
}

func (ph *ProblemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	err = ph.service.Delete(r.Context(), &datastruct.ProblemDeleteRequest{
		ProblemId: int32(id),
		ContestId: int32(contestId),
	})
	if err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, nil)
}
