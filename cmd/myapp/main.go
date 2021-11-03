package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	grpc_server "site/internal/grpc/server"
	"site/internal/http"
	"site/internal/store/inmemory"
	"syscall"

	"github.com/gorilla/sessions"
)

func main() {

	store := inmemory.NewDB()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	gsrv := grpc_server.NewServer(ctx, ":8081", store)
	go gsrv.Run()

	sessionStore := sessions.NewCookieStore([]byte("secret"))
	srv := http.NewServer(ctx, ":8080", store, sessionStore)

	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}