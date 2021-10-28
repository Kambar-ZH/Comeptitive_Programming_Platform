package main

import (
	"context"
	"log"
	grpc_server "site/grpc/server"
	"site/internal/http"
	"site/internal/store/inmemory"
)

func main() {

	store := inmemory.NewDB()

	gsrv := grpc_server.NewServer(context.Background(), ":8081", store)
	go gsrv.Run()

	srv := http.NewServer(context.Background(), ":8080", store)

	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()

}
