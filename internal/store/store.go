package store

import (
	"context"
	"site/internal/datastruct"
)

type Store interface {
	Connect(url string) error
	Close() error

	Users() UserRepository
	Submissions() SubmissionRepository
}

type UserRepository interface {
	All(ctx context.Context, offset, limit int) ([]*datastruct.User, error)
	ByEmail(ctx context.Context, email string) (*datastruct.User, error)
	ByHandle(ctx context.Context, handle string) (*datastruct.User, error)
	Create(ctx context.Context, user *datastruct.User) error
	Update(ctx context.Context, user *datastruct.User) error
	Delete(ctx context.Context, handle string) error
}

type SubmissionRepository interface {
	All(ctx context.Context, offset, limit int) ([]*datastruct.Submission, error)
	ByAuthorHandle(ctx context.Context, handle string) ([]*datastruct.Submission, error)
	ByContestId(ctx context.Context, contestId int) ([]*datastruct.Submission, error)
	ById(ctx context.Context, id int) (*datastruct.Submission, error)
	Create(ctx context.Context, submission *datastruct.Submission) error
	Update(ctx context.Context, submission *datastruct.Submission) error
	Delete(ctx context.Context, id int) error
}