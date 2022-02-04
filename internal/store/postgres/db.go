package postgres

import (
	"site/internal/logger"
	"site/internal/store"
	"time"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB

	users       store.UserRepository
	submissions store.SubmissionRepository
	contests    store.ContestRepository
	validators  store.ValidatorRepository
	testCases   store.TestCaseRepository
	problems    store.ProblemRepository
}

func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	start := time.Now()
	for conn.Ping() != nil {
		if start.After(start.Add(500 * time.Second)) {
			logger.Logger.Error("failed connect to db after 5 seconds")
			break
		}
	}

	if err := conn.Ping(); err != nil {
		return err
	}

	db.conn = conn
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
