package postgres

import (
	"context"
	"fmt"
	"log"
	"site/internal/datastruct"
	"testing"
)

const (
	urlAddress = "postgres://postgres:adminadmin@localhost:5432/codeforces"
)

func TestUserRepository(t *testing.T) {
	store := NewDB()
	if err := store.Connect(urlAddress); err != nil {
		log.Println(err)
		return
	}
	defer store.Close()
	ctx := context.Background()
	err := store.Users().Update(ctx, &datastruct.User{
		Handle:   "Kambar_Z",
		Email:    "zhamauov02@mail.ru",
		Country:  "Kazakhstan",
		City:     "Atyrau",
		Password: "password",
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	users, err := store.Users().All(ctx, 0, 10)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(users)
}
