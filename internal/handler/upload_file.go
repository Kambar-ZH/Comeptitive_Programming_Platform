package handler

import (
	"net/http"
	"site/internal/dto"
	"site/internal/services"
	"strconv"

	"github.com/go-chi/chi"
)

type UploadFileHandler struct {
	service services.UploadFileService
}

func NewUploadFileHandler(opts ...UploadFileHandlerOption) *UploadFileHandler {
	uf := &UploadFileHandler{}
	for _, v := range opts {
		v(uf)
	}
	return uf
}

func (uf UploadFileHandler) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, multipartFileHeader, err := r.FormFile("myFile")
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		problemId, err := strconv.Atoi(r.FormValue("problemId"))
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}
		contestId, err := strconv.Atoi(chi.URLParam(r, "contestId"))
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}
		submission, err := uf.service.UploadFile(r.Context(), &dto.UploadFileRequest{
			ProblemId: int32(problemId),
			ContestId: int32(contestId),
			File:      file,
			FileName:  multipartFileHeader.Filename,
		})
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		Respond(w, r, http.StatusAccepted, submission)
	}
}
