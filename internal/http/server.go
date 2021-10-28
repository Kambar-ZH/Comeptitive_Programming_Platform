package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"site/grpc/api"
	"site/internal/store"
	"site/test/compiler"
	"site/test/inmemory"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}

	store store.Store

	Address string
}

type ViewData struct {
	FileName    string
	ProblemName int
	Verdict     int
	FailedTest  int
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(inmemory.GetInstance().IndexHtml)
	if err != nil {
		fmt.Fprintf(w, "Error occured on loading home page")
		log.Println(err)
		return
	}
	data := ViewData{}
	tmpl.Execute(w, data)
}

func printFileInfo(handler *multipart.FileHeader) {
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)
}

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		return
	}
	defer file.Close()

	printFileInfo(handler)

	tempFile, err := ioutil.TempFile(inmemory.GetInstance().TempSolutions, "upload-*")
	if err != nil {
		log.Println("Error occured on creating temp file")
		log.Println(err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Error occured on reading file")
		log.Println(err)
		return
	}

	tempFile.Write(fileBytes)
	tempFile.Close()

	problemId, err := strconv.Atoi(r.FormValue("problemId"))
	if err != nil {
		log.Println("Error on parsing problem id")
		log.Println(err)
		return
	}
	fileName := strings.ReplaceAll(tempFile.Name(), "\\\\", "/")
	log.Println(fileName)

	////////////////////////////
	// 		Test solution 	  //
	////////////////////////////

	subResult := compiler.TestSolution(fileName, problemId)
	s.store.Submissions().Create(s.ctx, &api.Submission{
		Id:               15, // hardcoded
		ProblemId:        int32(problemId),
		AuthorId:         12,
		SubmissionResult: subResult,
	})

	////////////////////////////
	// Render upload template //
	////////////////////////////

	data := ViewData{
		FileName:    handler.Filename,
		ProblemName: problemId,
		Verdict:     int(subResult.Verdict),
		FailedTest:  int(subResult.FailedTest),
	}

	tmpl, err := template.ParseFiles(inmemory.GetInstance().UploadHtml)
	if err != nil {
		fmt.Fprintf(w, "Error occured on loading upload page")
		log.Println(err)
		return
	}
	tmpl.Execute(w, data)
}

func NewServer(ctx context.Context, addres string, store store.Store) *Server {
	return &Server{
		ctx:         ctx,
		Address:     addres,
		store:       store,
		idleConnsCh: make(chan struct{}),
	}
}

func (s *Server) SubmissionCrud(r chi.Router) chi.Router {
	r.Post("/submissions", func(w http.ResponseWriter, r *http.Request) {
		submission := new(api.Submission)
		if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		submissionResult, err := s.store.Submissions().Create(r.Context(), submission)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, submissionResult)
	})

	r.Get("/submissions", func(w http.ResponseWriter, r *http.Request) {
		submissions, err := s.store.Submissions().All(r.Context(), &api.Empty{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, submissions)
	})

	r.Get("/submissions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		submission, err := s.store.Submissions().ById(r.Context(), &api.SubmissionRequestId{Id: int32(id)})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

		render.JSON(w, r, submission)
	})

	r.Put("/submissions", func(w http.ResponseWriter, r *http.Request) {
		submission := new(api.Submission)
		if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		submissionResult, err := s.store.Submissions().Update(r.Context(), submission)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		render.JSON(w, r, submissionResult)
	})

	r.Delete("/submissions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.store.Submissions().Delete(r.Context(), &api.SubmissionRequestId{Id: int32(id)}); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}

func (s *Server) UserCrud(r chi.Router) chi.Router {
	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(api.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		userResult, err := s.store.Users().Create(r.Context(), user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, userResult)
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.Users().All(r.Context(), &api.Empty{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, users)
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, err := s.store.Users().ById(r.Context(), &api.UserRequestId{Id: int32(id)})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, user)
	})

	r.Put("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(api.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userResult, err := s.store.Users().Update(r.Context(), user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		render.JSON(w, r, userResult)
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err = s.store.Users().Delete(r.Context(), &api.UserRequestId{Id: int32(id)}); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	s.UserCrud(r)
	s.SubmissionCrud(r)

	r.HandleFunc("/upload", s.uploadFile)
	r.HandleFunc("/", s.homePage)

	return r
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}

	go s.ListenCtxForGT(srv)

	log.Printf("Serving on %v", srv.Addr)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // blocked until context not canceled

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Got error while shutting down %v", err)
		return
	}

	log.Println("Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}