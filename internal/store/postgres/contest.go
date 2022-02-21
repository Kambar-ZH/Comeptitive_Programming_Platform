package postgres

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/store"
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

func (c ContestRepository) FindAll(ctx context.Context, req *dto.ContestFindAllRequest) ([]*datastruct.Contest, error) {
	contests := make([]*datastruct.Contest, 0)
	if err := c.conn.Select(&contests,
		`SELECT 
		c.id "id",
       	c.start_date "start_date",
       	c.end_date "end_date",
       	c.phase "phase",
       	ct.name "name",
       	ct.description "description"
		FROM contests c
			JOIN contest_translation ct on c.id = ct.contest_id
		WHERE ct.language_code = $1
		ORDER BY start_date DESC
		OFFSET $2 
		LIMIT $3`,
		req.LanguageCode.String(), req.Offset, req.Limit); err != nil {
		return nil, err
	}
	return contests, nil
}

func (c ContestRepository) FindByTimeInterval(ctx context.Context, req *dto.ContestFindByTimeIntervalRequest) ([]*datastruct.Contest, error) {
	contests := make([]*datastruct.Contest, 0)
	if err := c.conn.Select(&contests,
		`SELECT 
		c.id "id",
       	c.start_date "start_date",
       	c.end_date "end_date",
       	c.phase "phase",
       	ct.name "name",
       	ct.description "description"
		FROM contests c
			JOIN contest_translation ct on c.id = ct.contest_id
		WHERE start_date BETWEEN $1 AND $2
			AND ct.language_code = $3`,
		req.StartTime, req.EndTime, req.LanguageCode.String()); err != nil {
		return nil, err
	}
	return contests, nil
}

func (c ContestRepository) GetById(ctx context.Context, req *dto.ContestGetByIdRequest) (*datastruct.Contest, error) {
	contest := new(datastruct.Contest)
	if err := c.conn.Get(contest,
		`SELECT 
		c.id "id",
       	c.start_date "start_date",
       	c.end_date "end_date",
       	c.phase "phase",
       	ct.name "name",
       	ct.description "description"
		FROM contests c
			JOIN contest_translation ct on c.id = ct.contest_id
		WHERE id = $1
			AND ct.language_code = $2`,
		req.ContestId, req.LanguageCode.String()); err != nil {
		return nil, err
	}
	return contest, nil
}

func (c ContestRepository) Create(ctx context.Context, contest *datastruct.Contest) error {
	tx, err := c.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, err = tx.Exec(
		`INSERT INTO 
			contests (start_date, end_date) 
			VALUES ($1, $2)`,
		contest.StartDate, contest.EndDate)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`INSERT INTO 
			contest_translation (contest_id, name, description) 
			VALUES (lastval(), $1, $2)`,
		contest.Name, contest.Description)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (c ContestRepository) Update(ctx context.Context, contest *datastruct.Contest) error {
	tx, err := c.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, err = tx.Exec(
		`UPDATE contests 
			SET start_date = $1, end_date = $2 
			WHERE id = $3`,
		contest.StartDate, contest.EndDate, contest.Id)
	_, err = tx.Exec(
		`UPDATE contest_translation 
			SET name = $1, description = $2 
			WHERE contest_id = $3`,
		contest.Name, contest.Description, contest.Id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (c ContestRepository) Delete(ctx context.Context, id int) error {
	tx, err := c.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, err = tx.Exec(
		`DELETE FROM contests 
			WHERE id = $1`, id)
	_, err = tx.Exec(
		`DELETE FROM contest_translation 
			WHERE contest_id = $1`, id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
