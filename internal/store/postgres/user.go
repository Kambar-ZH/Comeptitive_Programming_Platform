package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Users() store.UserRepository {
	if db.users == nil {
		db.users = NewUsersRepository(db.conn)
	}
	return db.users
}

type UserRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) store.UserRepository {
	return nil
}

func (u UserRepository) All(ctx context.Context, offset int, limit int) ([]*datastruct.User, error) {
	users := make([]*datastruct.User, 0)
	if err := u.conn.Select(&users, "SELECT * FROM users OFFSET = $1 LIMIT = $2", offset, limit); err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserRepository) ByEmail(ctx context.Context, email string) (*datastruct.User, error) {
	user := new(datastruct.User)
	if err := u.conn.Select(user, "SELECT * FROM users WHERE email = $1", email); err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserRepository) ByHandle(ctx context.Context, handle string) (*datastruct.User, error) {
	user := new(datastruct.User)
	if err := u.conn.Select(user, "SELECT * FROM users WHERE handle = $1", handle); err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserRepository) Create(ctx context.Context, user *datastruct.User) error {
	_, err := u.conn.Exec("INSERT INTO users(handle, email, country, city, password) VALUES ($1, $2, $3, $4, $5)",
		user.Handle, user.Email, user.Country, user.City, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u UserRepository) Update(ctx context.Context, user *datastruct.User) error {
	_, err := u.conn.Exec("UPDATE users SET(handle, email, country, city, password) VALUES ($1, $2, $3, $4, $5)",
		user.Handle, user.Email, user.Country, user.City, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u UserRepository) Delete(ctx context.Context, handle string) error {
	_, err := u.conn.Exec("DELETE FROM users WHERE handle = $1", handle)
	if err != nil {
		return err
	}
	return nil
}
