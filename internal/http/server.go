package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"site/internal/grpc/api"
	"site/internal/store"
	"site/test/compiler"
	"site/test/inmemory"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi"
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

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	log.Println(err)
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(inmemory.GetInstance().IndexHtml)
	if err != nil {
		s.error(w, r, http.StatusBadGateway, err)
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
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	printFileInfo(handler)

	tempFile, err := ioutil.TempFile(inmemory.GetInstance().TempSolutions, "upload-*")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	tempFile.Write(fileBytes)
	tempFile.Close()

	problemId, err := strconv.Atoi(r.FormValue("problemId"))
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
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
		s.error(w, r, http.StatusInternalServerError, err)
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
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		submissionResult, err := s.store.Submissions().Create(r.Context(), submission)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, submissionResult)
	})

	r.Get("/submissions", func(w http.ResponseWriter, r *http.Request) {
		submissions, err := s.store.Submissions().All(r.Context(), &api.Empty{})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, submissions)
	})

	r.Get("/submissions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		submission, err := s.store.Submissions().ById(r.Context(), &api.SubmissionRequestId{Id: int32(id)})
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, submission)
	})

	r.Put("/submissions", func(w http.ResponseWriter, r *http.Request) {
		submission := new(api.Submission)
		if err := json.NewDecoder(r.Body).Decode(submission); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		submissionResult, err := s.store.Submissions().Update(r.Context(), submission)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusAccepted, submissionResult)
	})

	r.Delete("/submissions/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err = s.store.Submissions().Delete(r.Context(), &api.SubmissionRequestId{Id: int32(id)}); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, "")
	})

	return r
}

func (s *Server) UserCrud(r chi.Router) chi.Router {
	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(api.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		userResult, err := s.store.Users().Create(r.Context(), user)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, userResult)
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.Users().All(r.Context(), &api.Empty{})
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, users)
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		user, err := s.store.Users().ById(r.Context(), &api.UserRequestId{Id: int32(id)})
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, user)
	})

	r.Put("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(api.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		userResult, err := s.store.Users().Update(r.Context(), user)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusAccepted, userResult)
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err = s.store.Users().Delete(r.Context(), &api.UserRequestId{Id: int32(id)}); err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		s.respond(w, r, http.StatusOK, "")
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

	log.Printf("serving on %v", srv.Addr)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // blocked until context not canceled

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("got error while shutting down %v", err)
		return
	}

	log.Println("proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}