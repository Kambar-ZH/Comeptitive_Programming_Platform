package http

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"site/internal/grpc/api"
	"site/internal/http/ioutils"
	"site/test/compiler"
	"site/test/inmemory"
	"strconv"
	"strings"
	"text/template"
)

func printFileInfo(handler *multipart.FileHeader) {
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)
}

func (s *Server) HomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetIndexHtml())
		if err != nil {
			ioutils.Error(w, r, http.StatusBadGateway, err)
			return
		}
		data := struct{}{}
		tmpl.Execute(w, data)
	}
}

func (s *Server) RegisterPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetRegisterHtml())
		if err != nil {
			ioutils.Error(w, r, http.StatusBadGateway, err)
			return
		}
		data := struct{}{}
		tmpl.Execute(w, data)
	}
}

func (s *Server) LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetLoginHtml())
		if err != nil {
			ioutils.Error(w, r, http.StatusBadGateway, err)
			return
		}
		data := struct{}{}
		tmpl.Execute(w, data)
	}
}

func (s *Server) RatingsPage() http.HandlerFunc {
	type RatingsViewData struct {
		Users   *api.UserList
		Country string
		City    string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		country := r.URL.Query().Get("country")
		if country == "" {
			country = "not configured"
		}
		city := r.URL.Query().Get("city")
		if city == "" {
			city = "not configured"
		}
		users, err := s.store.Users().All(s.ctx, &api.Pagination{})
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		data := RatingsViewData{
			Users:   users,
			Country: country,
			City:    city,
		}
		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetRatingsHtml())
		if err != nil {
			ioutils.Error(w, r, http.StatusBadGateway, err)
			return
		}
		tmpl.Execute(w, data)
	}
}

func (s *Server) UploadPage() http.HandlerFunc {
	type UploadViewData struct {
		FileName    string
		ProblemName int
		Verdict     int
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

		printFileInfo(handler)

		tempFile, err := ioutil.TempFile(inmemory.GetInstance().GetTempSolutions(), "upload-*")
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}

		tempFile.Write(fileBytes)
		tempFile.Close()

		problemId, err := strconv.Atoi(r.FormValue("problemId"))
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		fileName := strings.ReplaceAll(tempFile.Name(), "\\\\", "/")
		log.Println(fileName)

		subResult := compiler.TestSolution(fileName, problemId)
		s.store.Submissions().Create(s.ctx, &api.Submission{
			Id:               15, // hardcoded
			ProblemId:        int32(problemId),
			AuthorHandle:     "Kambar",
			SubmissionResult: subResult,
		})

		data := UploadViewData{
			FileName:    handler.Filename,
			ProblemName: problemId,
			Verdict:     int(subResult.Verdict),
			FailedTest:  int(subResult.FailedTest),
		}

		tmpl, err := template.ParseFiles(inmemory.GetInstance().GetUploadHtml())
		if err != nil {
			ioutils.Error(w, r, http.StatusInternalServerError, err)
			return
		}
		tmpl.Execute(w, data)
	}
}
