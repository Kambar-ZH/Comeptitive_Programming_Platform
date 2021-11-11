package http

import (
	"context"
	"log"
	"net/http"
	"site/internal/datastruct"
	"site/internal/http/handler"
	"site/internal/http/ioutils"
	"site/internal/middleware"
	"site/internal/services"
	"site/internal/store"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}

	store        store.Store
	sessionStore sessions.Store

	Address string
}

func NewServer(ctx context.Context, addres string, store store.Store, sessionStore sessions.Store) *Server {
	return &Server{
		ctx:          ctx,
		Address:      addres,
		store:        store,
		idleConnsCh:  make(chan struct{}),
		sessionStore: sessionStore,
	}
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	us := services.NewUserService(services.WithUserRepo(s.store.Users()))
	ss := services.NewSubmissionService(services.WithSubmissionRepo(s.store.Submissions()))
	ufs := services.NewUploadFileService(services.UploadFileWithSubmissionRepo(s.store.Submissions()))

	uh := handler.NewUserHandler(handler.WithUserService(us))
	sh := handler.NewSubmissionHandler(handler.WithSubmissionService(ss))
	ufh := handler.NewUploadFileHandler(handler.WithUploadFileService(ufs))

	s.PrepareUserRoutes(r, uh)
	s.PrepareSubmissionRoutes(r, sh)
	s.PrepareUploadFileRoutes(r, ufh)

	r.HandleFunc("/sessions", s.HandleSessionsCreate())
	r.Route("/private", func(r chi.Router) {
		r.Use(s.AuthenticateUser)
		r.HandleFunc("/whoami", s.handleWhoami())
	})
	return r
}

func (s *Server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ioutils.Respond(w, r, http.StatusOK, r.Context().Value(middleware.CtxKeyUser).(*datastruct.User))
	}
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
