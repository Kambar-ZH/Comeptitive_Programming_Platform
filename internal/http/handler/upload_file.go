package handler

import (
	"site/internal/services"
)

type UploadFileHandler struct {
	service services.UploadFileService
}

// func (ufh UploadFileHandler) UploadFile() http.HandlerFunc {
// 	type UploadViewData struct {
// 		FileName    string
// 		ProblemName int
// 		Verdict     string
// 		FailedTest  int
// 	}

// 	return func(w http.ResponseWriter, r *http.Request) {

// 		r.ParseMultipartForm(10 << 20)

// 		file, handler, err := r.FormFile("myFile")
// 		if err != nil {
// 			ioutils.Error(w, r, http.StatusInternalServerError, err)
// 			return
// 		}
// 		defer file.Close()

// 		tempFileName, err := ufh.service.SaveInmemory(file)
// 		if err != nil {
// 			ioutils.Error(w, r, http.StatusInternalServerError, err)
// 			return
// 		}

// 		problemId, err := strconv.Atoi(r.FormValue("problemId"))
// 		if err != nil {
// 			ioutils.Error(w, r, http.StatusInternalServerError, err)
// 			return
// 		}
// 		fileName := strings.ReplaceAll(tempFileName, "\\\\", "/")
// 		log.Println(fileName)

// 		subResult := ufh.service.TestSolution(fileName, problemId)
// 		s.store.Submissions().Create(s.ctx, &datastruct.Submission{
// 			Id:           15, // hardcoded
// 			ProblemId:    int32(problemId),
// 			AuthorHandle: "Kambar",
// 			Verdict:      string(subResult.Verdict),
// 		})

// 		data := UploadViewData{
// 			FileName:    handler.Filename,
// 			ProblemName: problemId,
// 			Verdict:     string(subResult.Verdict),
// 			FailedTest:  int(subResult.FailedTest),
// 		}

// 		ioutils.Respond(w, r, http.StatusAccepted, data)
// 	}
// }
