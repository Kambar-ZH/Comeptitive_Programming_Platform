package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"site/internal/http"
	"site/internal/message_broker/kafka"
	"site/internal/store/postgres"
	"syscall"

	"github.com/gorilla/sessions"
	lru "github.com/hashicorp/golang-lru"
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

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())

	go CatchTermination(cancel)
	
	sessionStore := sessions.NewCookieStore([]byte("secret"))

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer2")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(
		ctx,
		http.WithAddress(":8081"),
		http.WithStore(store),
		http.WithSessionStore(sessionStore),
		http.WithCache(cache),
		http.WithBroker(broker),
	)

	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}