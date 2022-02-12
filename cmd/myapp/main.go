package main

import (
	"context"
	"os"
	"os/signal"
	"site/internal/config"
	http "site/internal/http/rest"
	"site/internal/logger"
	"site/internal/message_broker/kafka"
	"site/internal/store/postgres"
	"syscall"

	"github.com/gorilla/sessions"
	lru "github.com/hashicorp/golang-lru"
)

// ctrl + / to make function

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

	brokers := []string{config.KafkaConn()}
	broker := kafka.NewBroker(brokers, cache, config.PeerName())
	if err := broker.Connect(ctx); err != nil {
		logger.Logger.Error(err.Error())
		os.Exit(1)
	}
	defer broker.Close()

	srv := http.NewServer(
		ctx,
		http.WithAddress(config.ServePort()),
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