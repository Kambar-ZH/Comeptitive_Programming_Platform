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
	"site/test"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	indexTemp       = "../../web/template/index.html"
	uploadTemp      = "../../web/template/upload.html"
	problemsStorage = "../../temp_solutions"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}

	store store.Store

	Address string
}

type ViewData struct {
	FileName    string
	ProblemName string
	Verdict     models.Verdict
	FailedTest  int
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(indexTemp)
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

func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return
	}
	defer file.Close()

	printFileInfo(handler)

	tempFile, err := ioutil.TempFile(problemsStorage, "upload-*")
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

	problemName := r.FormValue("problemName")
	fileName := strings.ReplaceAll(tempFile.Name(), "\\\\", "/")
	fmt.Println(fileName)

	verdict, failedTest := test.TestSolution(fileName, problemName)

	data := ViewData{
		FileName:    handler.Filename,
		ProblemName: problemName,
		Verdict:     verdict,
		FailedTest:  failedTest,
	}

	tmpl, err := template.ParseFiles(uploadTemp)
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

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()

	r.HandleFunc("/upload", uploadFile)
	r.HandleFunc("/", homePage)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			fmt.Fprintf(w, "Unknowr err: %v", err)
			return
		}

		s.store.Create(r.Context(), user)
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknowr err: %v", err)
			return
		}
		render.JSON(w, r, users)
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		user, err := s.store.ById(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, user)
	})

	r.Put("/users", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			fmt.Fprintf(w, "Unknowr err: %v", err)
			return
		}

		s.store.Create(r.Context(), user)
	})

	r.Delete("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Delete(r.Context(), id)
	})

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
