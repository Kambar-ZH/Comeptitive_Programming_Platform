package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"site/internal/grpc/api"
	"site/internal/http/ioutils"
	"site/internal/store"
	"strconv"

	"github.com/go-chi/chi"
)

type SubmissionHandler struct {
	Repository store.SubmissionRepository
}

func (sh *SubmissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	submission := new(api.Submission)
	if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
		ioutils.Error(w, r, http.StatusBadRequest, err)
		return
	}

	submissionResult, err := sh.Repository.Create(r.Context(), submission)
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	ioutils.Respond(w, r, http.StatusCreated, submissionResult)
}

func (sh *SubmissionHandler) All(w http.ResponseWriter, r *http.Request) {
	submissions, err := sh.Repository.All(r.Context(), &api.Pagination{})
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
	submission, err := sh.Repository.ById(r.Context(), &api.SubmissionRequestId{Id: int32(id)})
	if err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	if submission.ContestId != int32(contestId) {
		ioutils.Error(w, r, http.StatusMethodNotAllowed, fmt.Errorf("page preview not allowed"))
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

	submissionResult, err := sh.Repository.Update(r.Context(), submission)
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
	submission, err := sh.Repository.ById(r.Context(), &api.SubmissionRequestId{Id: int32(id)})
	if err != nil {
		ioutils.Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if submission.ContestId != int32(contestId) {
		ioutils.Error(w, r, http.StatusMethodNotAllowed, fmt.Errorf("page preview not allowed"))
		return
	}

	if _, err = sh.Repository.Delete(r.Context(), &api.SubmissionRequestId{Id: int32(id)}); err != nil {
		ioutils.Error(w, r, http.StatusNotFound, err)
		return
	}
	ioutils.Respond(w, r, http.StatusOK, "")
}
