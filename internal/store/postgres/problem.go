package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/logger"
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

func (p ProblemRepository) ProblemWithTags(ctx context.Context, problem *datastruct.Problem) (*datastruct.Problem, error) {
	tags := make([]*datastruct.Tag, 0)
	if err := p.conn.Select(&tags,
		`SELECT tags.* 
				FROM tags, problems_tags 
				WHERE problems_tags.problem_id = $1 
				AND tags.id = problems_tags.tag_id`,
		problem.Id); err != nil {
		return nil, err
	}
	problem.Tags = tags
	return problem, nil
}

func (p ProblemRepository) Problemset(ctx context.Context, req *dto.ProblemsetRequest) ([]*datastruct.Problem, error) {
	problems := make([]*datastruct.Problem, 0)
	if req.FilterTag != "" {
		if err := p.conn.Select(&problems,
			`SELECT DISTINCT problems.* 
				FROM problems, tags, problems_tags 
				WHERE problems_tags.tag_id = tags.id 
				AND tags.name = $1 
				AND problems.id = problems_tags.problem_id`,
			req.FilterTag); err != nil {
			return nil, err
		}
	} else {
		if err := p.conn.Select(&problems,
			`SELECT * 
				FROM problems 
				WHERE difficulty BETWEEN $1 AND $2 
				OFFSET $3 
				LIMIT $4`,
			req.MinDifficulty, req.MaxDifficulty, req.Offset, req.Limit); err != nil {
			return nil, err
		}
	}
	for i := range problems {
		var err error
		problems[i], err = p.ProblemWithTags(ctx, problems[i])
		if err != nil {
			logger.Logger.Error(err.Error())
		}
	}
	return problems, nil
}

func (p ProblemRepository) FindAll(ctx context.Context, req *dto.ProblemFindAllRequest) ([]*datastruct.Problem, error) {
	problems := make([]*datastruct.Problem, 0)
	if err := p.conn.Select(&problems,
		`SELECT * 
			FROM problems 
			WHERE contest_id = $1 
			OFFSET $2 
			LIMIT $3`,
		req.ContestId, req.Offset, req.Limit); err != nil {
		return nil, err
	}
	for i := range problems {
		var err error
		problems[i], err = p.ProblemWithTags(ctx, problems[i])
		if err != nil {
			logger.Logger.Error(err.Error())
		}
	}
	return problems, nil
}

func (p ProblemRepository) GetById(ctx context.Context, id int) (*datastruct.Problem, error) {
	problem := new(datastruct.Problem)
	if err := p.conn.Get(problem,
		`SELECT * 
			FROM problems 
			WHERE id = $1`,
		id); err != nil {
		return nil, err
	}
	var err error
	problem, err = p.ProblemWithTags(ctx, problem)
	if err != nil {
		return nil, err
	}
	return problem, nil
}

func (p ProblemRepository) Create(ctx context.Context, problem *datastruct.Problem) error {
	_, err := p.conn.Exec(
		`INSERT INTO problems(contest_id, index, name, statement, input, output, difficulty) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		problem.ContestId, problem.Index, problem.Name, problem.Statement, problem.Input, problem.Output, problem.Difficulty)
	if err != nil {
		return err
	}
	return nil
}

func (p ProblemRepository) Update(ctx context.Context, problem *datastruct.Problem) error {
	_, err := p.conn.Exec(
		`UPDATE problems 
			SET contest_id = $1, index = $2, name = $3, statement = $4, input = $5, output = $6, difficulty = $7 
			WHERE id = $8`,
		problem.ContestId, problem.Index, problem.Name, problem.Statement, problem.Input, problem.Output, problem.Difficulty, problem.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p ProblemRepository) Delete(ctx context.Context, id int) error {
	_, err := p.conn.Exec(
		`DELETE FROM problems 
			WHERE id = $1`,
		id)
	if err != nil {
		return err
	}
	return nil
}
