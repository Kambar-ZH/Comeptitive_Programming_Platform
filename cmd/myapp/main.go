package main

import (
	"context"
	"os"
	"os/signal"
	"site/internal/http"
	"site/internal/logger"
	"site/internal/message_broker/kafka"
	"site/internal/store/postgres"
	"site/internal/config"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	lru "github.com/hashicorp/golang-lru"

)

func main() {
	store := postgres.NewDB()
	if err := store.Connect(config.DSN()); err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())

	сatchTerminationfunc := func(cancel context.CancelFunc) {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		logger.Logger.Warn("caught termination signal")
		cancel()
	}

	go сatchTerminationfunc(cancel)

	sessionStore := sessions.NewCookieStore([]byte("secret"))

	brokers := []string{config.KAFKA_CONN()}
	broker := kafka.NewBroker(brokers, cache, config.PEER())
	for broker.Connect(ctx) != nil {
		time.Sleep(1 * time.Second)
		logger.Logger.Error("kafka cluster unreachable")
	}
	defer broker.Close()

	srv := http.NewServer(
		ctx,
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithSessionStore(sessionStore),
		http.WithCache(cache),
		http.WithBroker(broker),
	)

	if err := srv.Run(); err != nil {
		logger.Logger.Error(err.Error())
	}

	srv.WaitForGracefulTermination()
}
