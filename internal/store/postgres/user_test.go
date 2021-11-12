package postgres

import (
	"context"
	"fmt"
	"log"
	"site/internal/datastruct"
	"testing"
)


func TestUserRepository(t *testing.T) {
	store := NewDB()
	if err := store.Connect(urlAddress); err != nil {
		log.Println(err)
		return
	}
	defer store.Close()
	ctx := context.Background()
	
	// DELETE
	if err := store.Users().Delete(ctx, "Kambar_Z"); err != nil {
		panic(err)
	}
	// DELETE

	// CREATE
	err := store.Users().Create(ctx, &datastruct.User{
		Handle:            "Kambar_Z",
		Email:             "zhamauov02@mail.ru",
		Country:           "Kazakhstan",
		City:              "Atyrau",
		EncryptedPassword: "hash_password",
	})
	if err != nil {
		panic(err)
	}
	// CREATE

	// ALL
	users, err := store.Users().All(ctx, 0, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)
	// ALL

	// BYHANDLE
	user, err := store.Users().ByHandle(ctx, "Kambar_Z");
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	// BYHANDLE

	// UPDATE
	err = store.Users().Update(ctx, &datastruct.User{
		Handle:            "Kambar_Z",
		Email:             "yergeldi@mail.ru",
		Country:           "Kazakhstan",
		City:              "Uralsk",
		EncryptedPassword: "hash_password",
	})
	if err != nil {
		panic(err)
	}
	// UPDATE

	// BYEMAIL
	user, err = store.Users().ByEmail(ctx, "yergeldi@mail.ru");
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	// BYEMAIL
}