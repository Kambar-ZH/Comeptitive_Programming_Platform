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
	return db.contests
}

type ContestRepository struct {
	conn *sqlx.DB
}

func NewContestsRepository(conn *sqlx.DB) store.ContestRepository {
	return &ContestRepository{conn: conn}
}

func (c ContestRepository) All(ctx context.Context, query *datastruct.ContestAllRequest) ([]*datastruct.Contest, error) {
	contests := make([]*datastruct.Contest, 0)
	if err := c.conn.Select(&contests, "SELECT * FROM contests OFFSET $1 LIMIT $2", query.Offset, query.Limit); err != nil {
		return nil, err
	}
	return contests, nil
}

func (c ContestRepository) ById(ctx context.Context, id int) (*datastruct.Contest, error) {
	contest := new(datastruct.Contest)
	if err := c.conn.Get(contest, "SELECT * FROM contests WHERE id = $1", id); err != nil {
		return nil, err
	}
	return contest, nil
}

func (c ContestRepository) Create(ctx context.Context, contest *datastruct.Contest) error {
	_, err := c.conn.Exec("INSERT INTO contests(name, start_date, end_date, description) VALUES ($1, $2, $3, $4)",
		contest.Name, contest.StartDate, contest.EndDate, contest.Description)
	if err != nil {
		return err
	}
	return nil
}

func (c ContestRepository) Update(ctx context.Context, contest *datastruct.Contest) error {
	_, err := c.conn.Exec("UPDATE contests SET(name, start_date, end_date, description) VALUES ($1, $2, $3, $4)",
		contest.Name, contest.StartDate, contest.EndDate, contest.Description)
	if err != nil {
		return err
	}
	return nil
}

func (c ContestRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM contests WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
