package http

import (
	"context"
	"log"
	"net/http"
	"site/internal/http/handler"
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

	us := services.NewUserService(services.UserServiceWithStore(s.store))
	ss := services.NewSubmissionService(services.SubmissionServiceWithStore(s.store))
	ufs := services.NewUploadFileService(services.UploadFileServiceWithStore(s.store))
	as := services.NewAuthService(services.AuthServiceWithStore(s.store))

	uh := handler.NewUserHandler(handler.WithUserService(us))
	sh := handler.NewSubmissionHandler(handler.WithSubmissionService(ss))
	ufh := handler.NewUploadFileHandler(handler.WithUploadFileService(ufs))
	ah := handler.NewAuthHandler(handler.WithAuthService(as), handler.WithSessionStore(s.sessionStore))

	s.PrepareUserRoutes(r, uh)
	s.PrepareSubmissionRoutes(r, sh)
	s.PrepareUploadFileRoutes(r, ufh)

	r.HandleFunc("/sessions", ah.CreateSession())
	r.Route("/private", func(r chi.Router) {
		r.Use(ah.AuthenticateUser)
		r.HandleFunc("/whoami", ah.HandleWhoami())
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