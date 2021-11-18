package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"site/internal/datastruct"
	"site/internal/http/ioutils"
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
		ioutils.Error(w, r, http.StatusBadRequest, err)
		return
	}

	err := sh.service.Create(r.Context(), submission)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusCreated, nil)
}

func (sh *SubmissionHandler) All(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page := 1
	if (pageStr != "") {
		pageNum, err := strconv.Atoi(pageStr)
		if err != nil {
			ioutils.Error(w, r, http.StatusBadRequest, err)
			return
		}
		page = pageNum
	}
	submissions, err := sh.service.All(r.Context(), &datastruct.SubmissionQuery{Page: int32(page)})
	if err != nil {
		log.Println(err)
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, submissions)
}

func (sh *SubmissionHandler) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	submission, err := sh.service.ById(r.Context(), id)
	if err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, submission)
}

func (sh *SubmissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	submission := new(datastruct.Submission)
	if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	err := sh.service.Update(r.Context(), submission)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusAccepted, nil)
}

func (sh *SubmissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if err = sh.service.Delete(r.Context(), id); err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, "")
}