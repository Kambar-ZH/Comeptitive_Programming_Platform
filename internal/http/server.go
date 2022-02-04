package http

import (
	"context"
	"net/http"
	"site/internal/logger"
	messagebroker "site/internal/message_broker"
	"site/internal/store"
	"time"

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

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
	}

	go s.ListenCtxForGT(srv)

	logger.Logger.Sugar().Debugf("serving on %v", srv.Addr)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done() // blocked until context not canceled

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Logger.Sugar().Errorf("got error while shutting down %v", err.Error())
		return
	}

	logger.Logger.Debug("proccessed all idle connections")
	close(s.idleConnsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnsCh
}
