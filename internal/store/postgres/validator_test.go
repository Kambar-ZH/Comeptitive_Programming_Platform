package postgres_test

import (
	"context"
	"fmt"
	"site/internal/store/postgres"
	"site/internal/tools"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestValidators(t *testing.T) {
	conn, err := sqlx.Connect("pgx", "postgres://postgres:adminadmin@localhost:54320/codeforces")
	if err != nil {
		panic(err)
	}

	if err := conn.Ping(); err != nil {
		panic(err)
	}

	tester := postgres.NewValidatorRepository(conn)

	res, err := tester.GetByProblemId(context.Background(), 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)

	for _, tt := range res.TestCases {
		res, err := tools.MustExecuteFile(res.AuthorSolutionFilePath, tt)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
}
