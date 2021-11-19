package handler

import (
	"net/http"
	"site/internal/dto"
	"site/internal/services"
	"strconv"
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
	type UploadViewData struct {
		FileName    string
		ProblemName int
		Verdict     string
		FailedTest  int
	}

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("myFile")
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
		submission, err := uf.service.UploadFile(r.Context(), &dto.UploadFileRequest{
			ProblemId: problemId,
			File:      file,
		})
		if err != nil {
			Error(w, r, http.StatusInternalServerError, err)
			return
		}

		data := UploadViewData{
			FileName:    handler.Filename,
			ProblemName: problemId,
			Verdict:     string(submission.Verdict),
			FailedTest:  int(submission.FailedTest),
		}

		Respond(w, r, http.StatusAccepted, data)
	}
}
