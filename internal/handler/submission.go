package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"site/internal/datastruct"
	"site/internal/middleware"
	"site/internal/services"
	"strconv"

	"github.com/go-chi/chi"
)

type SubmissionHandler struct {
	service services.SubmissionService
}

func NewSubmissionHandler(opts ...SubmissionHandlerOption) *SubmissionHandler {
	sh := &SubmissionHandler{}
	for _, v := range opts {
		v(sh)
	}
	return sh
}

func (sh *SubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	submission := new(datastruct.Submission)
	if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
		Error(w, r, http.StatusBadRequest, err)
		return
	}

	err := sh.service.Create(r.Context(), submission)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, nil)
}

func (sh *SubmissionHandler) All(w http.ResponseWriter, r *http.Request) {
	query := &datastruct.SubmissionQuery{
		Page: r.Context().Value(middleware.CtxKeyPage).(int32),
		Filter: r.Context().Value(middleware.CtxKeyFilter).(string),
	}

	submissions, err := sh.service.All(r.Context(), query)
	if err != nil {
		log.Println(err)
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, submissions)
}

func (sh *SubmissionHandler) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	submission, err := sh.service.ById(r.Context(), id)
	if err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, submission)
}

func (sh *SubmissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	submission := new(datastruct.Submission)
	if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	err := sh.service.Update(r.Context(), submission)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusAccepted, nil)
}

func (sh *SubmissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if err = sh.service.Delete(r.Context(), id); err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, "")
}
