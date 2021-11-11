package handler

import (
	"log"
	"net/http"
	"site/internal/datastruct"
	"site/internal/http/ioutils"
	"site/internal/services"
	"strconv"
	"strings"
)

type UploadFileHandler struct {
	service services.UploadFileService
}

func NewUploadFileHandler(opts ...UploadFileHandlerOption) *UploadFileHandler {
	ufh := &UploadFileHandler{}
	for _, v := range(opts) {
		v(ufh)
	}
	return ufh
}

func (ufh UploadFileHandler) UploadFile() http.HandlerFunc {
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
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()

		tempFileName, err := ufh.service.SaveInmemory(file)
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		problemId, err := strconv.Atoi(r.FormValue("problemId"))
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		fileName := strings.ReplaceAll(tempFileName, "\\\\", "/")
		log.Println(fileName)

		subResult, err := ufh.service.TestSolution(fileName, problemId)
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		ufh.service.Create(r.Context(), &datastruct.Submission{
			Id:           15, // hardcoded
			ProblemId:    int32(problemId),
			Verdict:      string(subResult.Verdict),
		})

		data := UploadViewData{
			FileName:    handler.Filename,
			ProblemName: problemId,
			Verdict:     string(subResult.Verdict),
			FailedTest:  int(subResult.FailedTest),
		}

		ioutils.Respond(w, r, http.StatusAccepted, data)
	}
}