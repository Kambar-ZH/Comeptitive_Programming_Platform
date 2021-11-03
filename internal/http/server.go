package http

import (
	"context"
	"log"
	"net/http"
	"site/internal/grpc/api"
	"site/internal/http/handler"
	"site/internal/http/ioutils"
	"site/internal/store"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}
	logger      *logrus.Logger

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

	// r.Use(s.logRequest)

	uh := &handler.UserHandler{Repository: s.store.Users()}
	sh := &handler.SubmissionHandler{Repository: s.store.Submissions()}

	s.PrepareUserRoutes(r, uh)
	s.PrepareSubmissionRoutes(r, sh)

	s.PrepareGuestWebPageRoutes(r)
	
	r.HandleFunc("/sessions", s.HandleSessionsCreate())
	r.Route("/private", func(r chi.Router) {
		r.Use(s.AuthenticateUser)
		r.HandleFunc("/whoami", s.handleWhoami())
		s.PrepareAuthenticatedWebPageRoutes(r)
	})
	return r
}

func (s *Server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ioutils.Respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*api.User))
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
