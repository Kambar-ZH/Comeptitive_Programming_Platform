package http

import (
	"context"
	"log"
	"net/http"
	"site/internal/handler"
	messagebroker "site/internal/message_broker"
	"site/internal/services"
	"site/internal/store"
	"time"

	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	lru "github.com/hashicorp/golang-lru"
)

type Server struct {
	ctx         context.Context
	idleConnsCh chan struct{}

	store        store.Store
	sessionStore sessions.Store
	cache        *lru.TwoQueueCache
	broker       messagebroker.MessageBroker

	Address string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:         ctx,
		idleConnsCh: make(chan struct{}),
	}
	for _, v := range opts {
		v(srv)
	}
	return srv
}

func (s *Server) basicHandler() chi.Router {
	r := chi.NewRouter()
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	us := services.NewUserService(services.UserServiceWithStore(s.store))
	ss := services.NewSubmissionService(services.SubmissionServiceWithStore(s.store), services.SubmissionServiceWithBroker(s.broker), services.SubmissionServiceWithCache(s.cache))
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
