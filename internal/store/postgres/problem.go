package postgres

import (
	"context"
	"database/sql"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Problems() store.ProblemRepository {
	if db.problems == nil {
		db.problems = NewProblemRepository(db.conn)
	}
	return db.problems
}

type ProblemRepository struct {
	conn *sqlx.DB
}

func NewProblemRepository(conn *sqlx.DB) store.ProblemRepository {
	return &ProblemRepository{conn: conn}
}

func (p ProblemRepository) Problemset(ctx context.Context, req *dto.ProblemsetRequest) ([]*datastruct.Problem, error) {
	problems := make([]*datastruct.Problem, 0)
	if req.Filter != "" {
		if err := p.conn.Select(&problems, `
			SELECT
				p.id "id",
				p.contest_id "contest_id",
				p.difficulty "difficulty",
				p.points "points",
				p.index "index",
				pt.name "name",
				pt.input "input",
				pt.output "output",
				pt.statement "statement",
				pt.language_code "language_code"
			FROM problems p
				JOIN problem_translation pt on pt.problem_id = p.id 
				JOIN problems_tags pt2 on pt2.problem_id = p.id
				JOIN tags as t on pt2.tag_id = t.id
			WHERE t.name = $1
				AND pt.language_code = $2`,
			req.Filter, req.LanguageCode); err != nil {
			return nil, err
		}
	} else {
		if err := p.conn.Select(&problems, `
			SELECT
				p.id "id",
				p.contest_id "contest_id",
				p.difficulty "difficulty",
				p.points "points",
				p.index "index",
				pt.name "name",
				pt.input "input",
				pt.output "output",
				pt.statement "statement",
				pt.language_code "language_code"
			FROM problems p
				JOIN problem_translation pt on p.id = pt.problem_id
			WHERE p.difficulty BETWEEN $1 AND $2 
				AND pt.language_code = $3
			OFFSET $4
			LIMIT $5`,
			req.MinDifficulty, req.MaxDifficulty, req.LanguageCode.String(), req.Offset, req.Limit); err != nil {
			return nil, err
		}
	}
	return problems, nil
}

func (p ProblemRepository) FindAll(ctx context.Context, req *dto.ProblemFindAllRequest) ([]*datastruct.Problem, error) {
	problems := make([]*datastruct.Problem, 0)
	err := p.conn.Select(&problems, `
		SELECT
			p.id "id",
			p.contest_id "contest_id",
			p.difficulty "difficulty",
			p.points "points",
			p.index "index",
			pt.name "name",
			pt.input "input",
			pt.output "output",
			pt.statement "statement",
		    pt.language_code "language_code"
		FROM problems p
			JOIN problem_translation pt on p.id = pt.problem_id
		WHERE p.contest_id = $1 
			AND pt.language_code = $2`,
		req.ContestId, req.LanguageCode.String())
	if err != nil {
		return nil, err
	}
	return problems, nil
}

func (p ProblemRepository) GetById(ctx context.Context, req *dto.ProblemGetByIdRequest) (*datastruct.Problem, error) {
	problem := new(datastruct.Problem)
	err := p.conn.Get(problem, `
		SELECT 
			p.id "id",
       		p.contest_id "contest_id",
       		p.difficulty "difficulty",
       		p.points "points",
       		p.index "index",
			pt.name "name",
			pt.input "input",
       		pt.output "output",
       		pt.statement "statement",
		       		pt.language_code "language_code"
		FROM problems p
			JOIN problem_translation pt on p.id = pt.problem_id
		WHERE p.id = $1
			AND pt.language_code = $2`,
		req.ProblemId, req.LanguageCode.String())
	if err != nil {
		return nil, err
	}
	return problem, nil
}

func (p ProblemRepository) Create(ctx context.Context, problem *datastruct.Problem) error {
	tx, err := p.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO problems(contest_id, index, difficulty, points) 
			VALUES ($1, $2, $3, $4)`,
		problem.ContestId, problem.Index, problem.Difficulty, problem.Points)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO problem_translation(problem_id, name, statement, input, output) 
			VALUES (lastval(), $1, $2, $3, $4)`,
		problem.Name, problem.Statement, problem.Input, problem.Output)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p ProblemRepository) Update(ctx context.Context, problem *datastruct.Problem) error {
	tx, err := p.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, err = tx.Exec(`
		UPDATE problems 
			SET contest_id = $1, index = $2, difficulty = $3, points = $4
		WHERE id = $4`,
		problem.ContestId, problem.Index, problem.Difficulty, problem.Id, problem.Points)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE problem_translation 
			SET name = $1, statement = $2, input = $3, output = $4
		WHERE problem_id = $5`,
		problem.Name, problem.Statement, problem.Input, problem.Output, problem.Id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p ProblemRepository) Delete(ctx context.Context, id int) error {
	tx, err := p.conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	_, err = tx.Exec(`
		DELETE FROM problems 
		WHERE id = $1`,
		id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		DELETE FROM problem_translation
		WHERE problem_id = $1`,
		id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
