package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Contests() store.ContestRepository {
	if db.contests == nil {
		db.contests = NewContestsRepository(db.conn)
	}
}

type ContestRepository struct {
	conn *sqlx.DB
}

func NewContestsRepository(conn *sqlx.DB) ContestRepository {
	return ContestRepository{conn: conn}
}

func (c ContestRepository) All(ctx context.Context, query *datastruct.ContestQuery) ([]*datastruct.Contest, error) {
	contests := make([]*datastruct.Contest, 0)
	if err := c.conn.Select(contests, "SELECT * FROM contests OFFSET $1 LIMIT $2", query.Offset, query.Limit); err != nil {
		return nil, err
	}
	return contests, nil
}

func (c ContestRepository) ById(ctx context.Context, id int) (*datastruct.Contest, error) {
	contest := new(datastruct.Contest)
	if err := c.conn.Select(contest, "SELECT * FROM contests WHERE id = $1", id); err != nil {
		return nil, err
	}
	return contest, nil
}

func (c ContestRepository) Create(ctx context.Context, contest *datastruct.Contest) error {
	_, err := c.conn.Exec("INSERT INTO contests(")
}

func (c ContestRepository) Update(ctx context.Context, contest *datastruct.Contest) error {
	panic("not implemented") // TODO: Implement
}

func (c ContestRepository) Delete(ctx context.Context, id int) error {
	panic("not implemented") // TODO: Implement
}
