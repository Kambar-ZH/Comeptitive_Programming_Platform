package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"site/internal/datastruct"
	"site/internal/dto"
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
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}

	err = sh.service.Create(r.Context(), &dto.SubmissionCreateRequest{
		Submission: submission,
		ContestId:  int32(contestId),
	})
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusCreated, nil)
}

func (sh *SubmissionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	req := &dto.SubmissionFindAllRequest{
		Page:             r.Context().Value(middleware.CtxKeyPage).(int32),
		FilterUserHandle: r.Context().Value(middleware.CtxKeyFilter).(string),
		ContestId:        int32(contestId),
	}

	submissions, err := sh.service.All(r.Context(), req)
	if err != nil {
		log.Println(err)
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	Respond(w, r, http.StatusOK, submissions)
}

func (sh *SubmissionHandler) GetById(w http.ResponseWriter, r *http.Request) {
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
	submission, err := sh.service.ById(r.Context(), &dto.SubmissionGetByIdRequest{
		SubmissionId: int32(id),
		ContestId:    int32(contestId),
	})
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
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if err := sh.service.Update(r.Context(), &dto.SubmissionUpdateRequest{
		Submission: submission,
		ContestId:  int32(contestId),
	}); err != nil {
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
	contestIdStr := chi.URLParam(r, "contestId")
	contestId, err := strconv.Atoi(contestIdStr)
	if err != nil {
		Error(w, r, http.StatusInternalServerError, err)
		return
	}
	if err = sh.service.Delete(r.Context(), &dto.SubmissionDeleteRequest{
		SubmissionId: int32(id),
		ContestId:    int32(contestId),
	}); err != nil {
		Error(w, r, http.StatusNotFound, err)
		return
	}
	Respond(w, r, http.StatusOK, nil)
}
