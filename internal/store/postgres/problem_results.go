package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) ProblemResults() store.ProblemResultsRepository {
	if db.problemResults == nil {
		db.problemResults = NewProblemResultsRepository(db.conn)
	}
	return db.problemResults
}

type ProblemResultsRepository struct {
	conn *sqlx.DB
}

func NewProblemResultsRepository(conn *sqlx.DB) store.ProblemResultsRepository {
	return &ProblemResultsRepository{conn: conn}
}

func (p ProblemResultsRepository) ProblemResultsWithSubmissions(ctx context.Context, problemResult *datastruct.ProblemResult) (*datastruct.ProblemResult, error) {
	submissions := make([]*datastruct.Submission, 0)
	if err := p.conn.Select(&submissions,
		`SELECT * 
			FROM submissions 
			WHERE contest_id = $1 
			AND user_id = $2 
			AND problem_id = $3`,
		problemResult.ContestId, problemResult.UserId, problemResult.ProblemId); err != nil {
		return nil, err
	}
	problemResult.Submissions = submissions
	return problemResult, nil
}

func (p ProblemResultsRepository) GetByProblemId(ctx context.Context, req *dto.ProblemResultGetByProblemIdRequest) (*datastruct.ProblemResult, error) {
	problemResult := new(datastruct.ProblemResult)
	if err := p.conn.Get(problemResult,
		`SELECT * 
			FROM problem_results 
			WHERE user_id = $1
			AND problem_id = $2
			AND contest_id = $3`,
		req.UserId, req.ProblemId, req.ContestId); err != nil {
		return nil, err
	}

	problemResult, err := p.ProblemResultsWithSubmissions(ctx, problemResult)
	if err != nil {
		return nil, err
	}
	return problemResult, nil
}

func (p ProblemResultsRepository) Update(ctx context.Context, problemResults *datastruct.ProblemResult) error {
	_, err := p.conn.Exec(
		`UPDATE problem_results 
			SET user_id = $1, contest_id = $2, problem_id = $3, penalty = $4, points = $5, last_successful_submission_time = $6
			WHERE user_id = $1
			AND contest_id = $2
			AND problem_id = $3`,
		problemResults.UserId, problemResults.ContestId, problemResults.ProblemId, problemResults.Penalty, problemResults.Points, problemResults.LastSuccessfulSubmissionTime)
	if err != nil {
		return err
	}
	return nil
}

func (p ProblemResultsRepository) Create(ctx context.Context, problemResults *datastruct.ProblemResult) error {
	_, err := p.conn.Exec(
		`INSERT INTO 
			problem_results (user_id, contest_id, problem_id, penalty, points, last_successful_submission_time)
			VALUES($1, $2, $3, $4, $5, $6)`,
		problemResults.UserId, problemResults.ContestId, problemResults.ProblemId, problemResults.Penalty, problemResults.Points, problemResults.LastSuccessfulSubmissionTime)
	if err != nil {
		return err
	}
	return nil
}
