package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"site/internal/models"
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
	Verdict     models.Verdict
	FailedTest  int
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(inmemory.GetInstance().IndexHtml)
	if err != nil {
		fmt.Fprintf(w, "Error occured on loading home page")
		fmt.Println(err)
		return
	}
	data := ViewData{}
	tmpl.Execute(w, data)
}

func printFileInfo(handler *multipart.FileHeader) {
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
}

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return
	}
	defer file.Close()

	printFileInfo(handler)

	tempFile, err := ioutil.TempFile(inmemory.GetInstance().TempSolutions, "upload-*")
	if err != nil {
		fmt.Println("Error occured on creating temp file")
		fmt.Println(err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error occured on reading file")
		fmt.Println(err)
		return
	}

	tempFile.Write(fileBytes)
	tempFile.Close()

	problemId, err := strconv.Atoi(r.FormValue("problemId"))
	if err != nil {
		fmt.Println("Error on parsing problem id")
		fmt.Println(err)
		return
	}
	fileName := strings.ReplaceAll(tempFile.Name(), "\\\\", "/")
	fmt.Println(fileName)

	////////////////////////////
	// 		Test solution 	  //
	////////////////////////////

	subResult := compiler.TestSolution(fileName, problemId)
	s.store.Submissions().Create(s.ctx, &models.Submission{
		Id:               15, // hardcoded
		ProblemId:        problemId,
		AuthorId:         12,
		SubmissionResult: subResult,
	})

	////////////////////////////
	// Render upload template //
	////////////////////////////

	data := ViewData{
		FileName:    handler.Filename,
		ProblemName: problemId,
		Verdict:     subResult.Verdict,
		FailedTest:  subResult.FailedTest,
	}

	tmpl, err := template.ParseFiles(inmemory.GetInstance().UploadHtml)
	if err != nil {
		fmt.Fprintf(w, "Error occured on loading upload page")
		fmt.Println(err)
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
		submission := new(models.Submission)
		if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := s.store.Submissions().Create(r.Context(), submission); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/submissions", func(w http.ResponseWriter, r *http.Request) {
		submissions, err := s.store.Submissions().All(r.Context())
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

		submission, err := s.store.Submissions().ById(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

		render.JSON(w, r, submission)
	})

	r.Put("/submissions", func(w http.ResponseWriter, r *http.Request) {
		submission := new(models.Submission)
		if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.store.Submissions().Update(r.Context(), submission); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})

	r.Delete("/submissions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = s.store.Submissions().Delete(r.Context(), id); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}

func (s *Server) UserCrud(r chi.Router) chi.Router {
	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := s.store.Users().Create(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.Users().All(r.Context())
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

		user, err := s.store.Users().ById(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, user)
	})

	r.Put("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.store.Users().Update(r.Context(), user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = s.store.Users().Delete(r.Context(), id); err != nil {
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

	fmt.Println("Server started work")
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // blocked until context not canceled

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("Got error while shutting down %v", err)
		return
	}

	fmt.Println("Proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}
