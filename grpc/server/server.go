package server

import (
	"context"
	"log"
	"net"
	"site/grpc/api"
	"site/internal/store"

	"google.golang.org/grpc"
)

type Server struct {
	ctx     context.Context
	store   store.Store
	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:     ctx,
		Address: address,
		store:   store,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		log.Fatalf("cannot listen to %s: %v", s.Address, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()

	api.RegisterUserRepositoryServer(grpcServer, s.store.Users())

	log.Printf("Serving on %v", listener.Addr())
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve on %v: %v", listener.Addr(), err)
	}
}