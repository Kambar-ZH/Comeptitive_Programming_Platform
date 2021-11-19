package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"site/internal/http"
	"site/internal/store/cache"
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

	cache := cache.NewRedisCache(
		"localhost:6379",
		0,
		10,
	)
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	sessionStore := sessions.NewCookieStore([]byte("secret"))

	srv := http.NewServer(
		ctx,
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithSessionStore(sessionStore),
		http.WithCache(cache),
	)

	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
