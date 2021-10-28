package main

import (
	"context"
	"site/grpc/api"
	"log"

	"google.golang.org/grpc"
)

const (
	port = ":8081"
)

func main() {
	ctx := context.Background()

	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", port, err)
	}

	userRepositoryClient := api.NewUserRepositoryClient(conn)

	newUser := &api.User{
		Id:        1,
		Handle:    "Kambar",
		Country:   "Kazakhstan",
		City:      "Atyrau",
		Rating:    1703,
		MaxRating: 1703,
		Avatar:    "url-link",
	}

	createResponseUser, err := userRepositoryClient.Create(ctx, newUser)

	users, err := userRepositoryClient.All(ctx, &api.Empty{})
	if err != nil {
		log.Fatalf("could not get users: %v", err)
	}
	log.Printf("got list of users: %v", users.Users)

	validId, invalidId := int32(1), int32(3)
	user, err := userRepositoryClient.ById(ctx, &api.UserRequestId{Id: validId})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got user with id %d: %v", validId, user)

	_, err = userRepositoryClient.ById(ctx, &api.UserRequestId{Id: invalidId})
	if err != nil {
		log.Printf("got error: %v", err)
	}

	newUser = &api.User{
		Id:        2,
		Handle:    "Yergeldi",
		Country:   "Kazakhstan",
		City:      "Ust Kamenogorsk",
		Rating:    2500,
		MaxRating: 3000,
		Avatar:    "url-link",
	}

	createResponseUser, err = userRepositoryClient.Create(ctx, newUser)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user successfully createed to db: %v", createResponseUser)

	updatedUser := &api.User{
		Id:        1,
		Handle:    "Kambar_Z",
		Country:   "Kazakhstan",
		City:      "Atyrau",
		Rating:    1800,
		MaxRating: 1800,
		Avatar:    "url-link",
	}
	updateResponseUser, err := userRepositoryClient.Update(ctx, updatedUser)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user successfully updated: %v", updateResponseUser)

	_, err = userRepositoryClient.Delete(ctx, &api.UserRequestId{Id: 1})
	if err != nil {
		log.Fatal(err)
	}
}
