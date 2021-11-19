package http

import (
	"site/internal/store"
	"site/internal/cache"

	"github.com/gorilla/sessions"
)

type ServerOption func(srv *Server)

func WithAddress(address string) ServerOption {
	return func(srv *Server) {
		srv.Address = address
	}
}

func WithStore(store store.Store) ServerOption {
	return func(srv *Server) {
		srv.store = store
	}
}

func WithSessionStore(sessionStore sessions.Store) ServerOption {
	return func(srv *Server) {
		srv.sessionStore = sessionStore
	}
}

func WithCache(cache cache.Cache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}
