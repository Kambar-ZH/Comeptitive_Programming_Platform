package main

import (
	"fmt"
	"site/internal/datastruct"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	urlAddress := "postgres://postgres:adminadmin@0.0.0.0:54320/codeforces"
	conn, err := sqlx.Connect("pgx", urlAddress)
	if err != nil {
		panic(err)
	}
	users := make([]*datastruct.Contest, 0)
	if err := conn.Select(&users, `SELECT * 
	FROM contests`); err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}
