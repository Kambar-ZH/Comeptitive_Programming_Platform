package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"site/internal/grpc/api"
	"site/internal/http/ioutils"
	"site/internal/services"
	"strconv"

	"github.com/go-chi/chi"
)

type SubmissionHandler struct {
	service services.SubmissionService
}

func NewSubmissionService(opts ...SubmissionHandlerOption) *SubmissionHandler {
	sh := &SubmissionHandler{}
	for _, v := range opts {
		v(sh)
	}
	return sh
}

func (sh *SubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	submission := new(api.Submission)
	if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
		ioutils.Error(w, r, http.StatusBadRequest, err)
		return
	}

	submissionResult, err := sh.service.Create(r.Context(), submission)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusCreated, submissionResult)
}

func (sh *SubmissionHandler) All(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		ioutils.Error(w, r, http.StatusBadRequest, err)
	}
	submissions, err := sh.service.All(r.Context(), &api.AllSubmissionsRequest{Page: int32(page)})
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
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		ioutils.Error(w, r, http.StatusBadRequest, err)
	}
	submission, err := sh.service.ById(r.Context(), &api.SubmissionByIdRequest{Id: int32(id)})
	if err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, submission)
}

func (sh *SubmissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		ioutils.Error(w, r, http.StatusBadRequest, err)
	}
	submission := new(api.Submission)
	if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if submission.ContestId != int32(contestId) {
		ioutils.Error(w, r, http.StatusMethodNotAllowed, fmt.Errorf("page preview not allowed"))
		return
	}

	submissionResult, err := sh.service.Update(r.Context(), submission)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusAccepted, submissionResult)
}

func (sh *SubmissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		ioutils.Error(w, r, http.StatusBadRequest, err)
	}
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	submission, err := sh.service.ById(r.Context(), int32(id))
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if submission.ContestId != int32(contestId) {
		ioutils.Error(w, r, http.StatusMethodNotAllowed, fmt.Errorf("page preview not allowed"))
		return
	}

	if _, err = sh.service.Delete(r.Context(), int32(id)); err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, "")
}
