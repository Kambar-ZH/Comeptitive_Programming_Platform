package postgres

import (
	_ "github.com/jackc/pgx/stdlib"
	"site/internal/store"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB

	users       store.UserRepository
	submissions store.SubmissionRepository
}

func NewDB() store.Store {
	return &DB{}
}

func (db *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
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
