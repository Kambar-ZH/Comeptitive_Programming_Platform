package postgres

import (
	"context"
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
	return nil
}

func (u SubmissionRepository) All(ctx context.Context, offset int, limit int) ([]*datastruct.Submission, error) {
	submissions := make([]*datastruct.Submission, 0)
	if err := u.conn.Select(&submissions, "SELECT * FROM submissions OFFSET = $1 LIMIT = $2", offset, limit); err != nil {
		return nil, err
	}
	return submissions, nil
}

// All(ctx context.Context, offset, limit int) ([]*datastruct.Submission, error)
//     ByAuthorHandle(ctx context.Context, handle string) ([]*datastruct.Submission, error)
//     ByContestId(ctx context.Context, contestId int) ([]*datastruct.Submission, error)
//     ById(ctx context.Context, id int) (*datastruct.Submission, error)
//     Create(ctx context.Context, submission *datastruct.Submission) error
//     Update(ctx context.Context, submission *datastruct.Submission) error
//     Delete(ctx context.Context, id int) error

func (u SubmissionRepository) ById(ctx context.Context, id int) (*datastruct.Submission, error) {
	submission := new(datastruct.Submission)
	if err := u.conn.Select(submission, "SELECT * FROM submissions WHERE id = $1", id); err != nil {
		return nil, err
	}
	return submission, nil
}

func (u SubmissionRepository) ByContestId(ctx context.Context, contestId int) (*datastruct.Submission, error) {
	submission := new(datastruct.Submission)
	if err := u.conn.Select(submission, "SELECT * FROM submissions WHERE contest_id = $1", contestId); err != nil {
		return nil, err
	}
	return submission, nil
}

func (u SubmissionRepository) ByAuthorHandle(ctx context.Context, handle string) (*datastruct.Submission, error) {
	submission := new(datastruct.Submission)
	if err := u.conn.Select(submission, "SELECT * FROM submissions WHERE author_handle = $1", handle); err != nil {
		return nil, err
	}
	return submission, nil
}

func (u SubmissionRepository) Create(ctx context.Context, submission *datastruct.Submission) error {
	// TODO Date conversion
	_, err := u.conn.Exec("INSERT INTO submissions(contest_id, problem_id, author_handle, verdict, failed_test) VALUES ($1, $2, $3, $4, $5)", 
		submission.ContestId, submission.ProblemId, submission.AuthorHandle, submission.Verdict, submission.FailedTest)
	if err != nil {
		return err
	}
	return nil
}

func (u SubmissionRepository) Update(ctx context.Context, submission *datastruct.Submission) error {
	// TODO Date conversion
	_, err := u.conn.Exec("UPDATE submissions SET(contest_id, problem_id, author_handle, verdict, failed_test) VALUES ($1, $2, $3, $4, $5)", 
		submission.ContestId, submission.ProblemId, submission.AuthorHandle, submission.Verdict, submission.FailedTest)
	if err != nil {
		return err
	}
	return nil	
}

func (u SubmissionRepository) Delete(ctx context.Context, id int) error {
	_, err := u.conn.Exec("DELETE FROM submissions WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}