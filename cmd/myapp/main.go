package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"site/internal/http"
	"site/internal/store/postgres"
	"syscall"

	"github.com/gorilla/sessions"
)

const (
	urlAddress = "postgres://postgres:adminadmin@localhost:5432/codeforces"
)

func main() {
	store := postgres.NewDB()
	if err := store.Connect(urlAddress); err != nil {
		panic(err)
	}
	defer store.Close()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	sessionStore := sessions.NewCookieStore([]byte("secret"))
	srv := http.NewServer(ctx, ":8080", store, sessionStore)

	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
