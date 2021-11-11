package server

// import (
// 	"context"
// 	"log"
// 	"net"
// 	"site/internal/grpc/api"
// 	"site/internal/store"

// 	"google.golang.org/grpc"
// )

// type Server struct {
// 	ctx         context.Context
// 	store       store.Store
// 	Address     string
// 	idleConnsCh chan struct{}
// }

// func NewServer(ctx context.Context, address string, store store.Store) *Server {
// 	return &Server{
// 		ctx:         ctx,
// 		Address:     address,
// 		store:       store,
// 		idleConnsCh: make(chan struct{}),
// 	}
// }

// func (s *Server) Run() {
// 	listener, err := net.Listen("tcp", s.Address)
// 	if err != nil {
// 		log.Fatalf("cannot listen to %s: %v", s.Address, err)
// 	}
// 	defer listener.Close()

// 	grpcServer := grpc.NewServer()

// 	api.RegisterUserRepositoryServer(grpcServer, s.store.Users())
// 	api.RegisterSubmissionRepositoryServer(grpcServer, s.store.Submissions())

// 	go s.ListenCtxForGT(grpcServer)

// 	log.Printf("serving on %v", listener.Addr())
// 	if err := grpcServer.Serve(listener); err != nil {
// 		log.Fatalf("failed to serve on %v: %v", listener.Addr(), err)
// 	}
// 	log.Println("grpc: server closed")
// }

// func (s *Server) ListenCtxForGT(srv *grpc.Server) {
// 	<-s.ctx.Done() // blocked until context not canceled

// 	srv.GracefulStop()

// 	log.Println("proccessed all idle connections")
// 	close(s.idleConnsCh)
// }

// func (s *Server) WaitForGracefulTermination() {
// 	<-s.idleConnsCh
// }
