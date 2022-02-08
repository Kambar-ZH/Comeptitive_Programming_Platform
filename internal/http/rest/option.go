package http

import (
	messagebroker "site/internal/message_broker"
	"site/internal/store"

	"github.com/gorilla/sessions"
	lru "github.com/hashicorp/golang-lru"
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

func WithCache(cache *lru.TwoQueueCache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}

func WithBroker(broker messagebroker.MessageBroker) ServerOption {
	return func(srv *Server) {
		srv.broker = broker
	}
}
