package postgres

import (
	"context"
	"fmt"
	"site/internal/datastruct"
	"site/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Submissions() store.SubmissionRepository {
	if db.submissions == nil {
		db.submissions = NewSubmissionRepository(db.conn)
	}
	return db.submissions
}

type SubmissionRepository struct {
	conn *sqlx.DB
}

func NewSubmissionRepository(conn *sqlx.DB) store.SubmissionRepository {
	return &SubmissionRepository{conn: conn}
}

func (s SubmissionRepository) All(ctx context.Context, req *datastruct.SubmissionAllRequest) ([]*datastruct.Submission, error) {
	submissions := make([]*datastruct.Submission, 0)
	basicQuery := "SELECT * FROM submissions"
	if req.Filter != "" {
		basicQuery = fmt.Sprintf("%s WHERE author_handle ILIKE $1 and contest_id = $2 OFFSET $3 LIMIT $4", basicQuery)

		if err := s.conn.Select(&submissions, basicQuery, req.ContestId, "%"+req.Filter+"%", req.Offset, req.Limit); err != nil {
			return nil, err
		}
		return submissions, nil
	}
	if err := s.conn.Select(&submissions, "SELECT * FROM submissions WHERE contest_id = $1 OFFSET $2 LIMIT $3", req.ContestId, req.Offset, req.Limit); err != nil {
		return nil, err
	}
	return submissions, nil
}

func (s SubmissionRepository) ById(ctx context.Context, id int) (*datastruct.Submission, error) {
	submission := new(datastruct.Submission)
	if err := s.conn.Get(submission, "SELECT * FROM submissions WHERE id = $1", id); err != nil {
		return nil, err
	}
	return submission, nil
}

func (s SubmissionRepository) Create(ctx context.Context, submission *datastruct.Submission) error {
	// TODO: Date conversion
	_, err := s.conn.Exec("INSERT INTO submissions(contest_id, problem_id, author_handle, verdict, failed_test) VALUES ($1, $2, $3, $4, $5)",
		submission.ContestId, submission.ProblemId, submission.AuthorHandle, submission.Verdict, submission.FailedTest)
	if err != nil {
		return err
	}
	return nil
}

func (s SubmissionRepository) Update(ctx context.Context, submission *datastruct.Submission) error {
	// TODO: Date conversion
	_, err := s.conn.Exec("UPDATE submissions SET(contest_id, problem_id, author_handle, verdict, failed_test) VALUES ($1, $2, $3, $4, $5)",
		submission.ContestId, submission.ProblemId, submission.AuthorHandle, submission.Verdict, submission.FailedTest)
	if err != nil {
		return err
	}
	return nil
}

func (s SubmissionRepository) Delete(ctx context.Context, id int) error {
	_, err := s.conn.Exec("DELETE FROM submissions WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
