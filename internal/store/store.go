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
	Validators() ValidatorRepository
	TestCases() TestCaseRepository
	Contests() ContestRepository
}

type UserRepository interface {
	All(ctx context.Context, query *datastruct.UserQuery) ([]*datastruct.User, error)
	ByEmail(ctx context.Context, email string) (*datastruct.User, error)
	ByHandle(ctx context.Context, handle string) (*datastruct.User, error)
	Create(ctx context.Context, user *datastruct.User) error
	Update(ctx context.Context, user *datastruct.User) error
	Delete(ctx context.Context, handle string) error
}

type SubmissionRepository interface {
	All(ctx context.Context, query *datastruct.SubmissionAllRequest) ([]*datastruct.Submission, error)
	ById(ctx context.Context, id int) (*datastruct.Submission, error)
	Create(ctx context.Context, submission *datastruct.Submission) error
	Update(ctx context.Context, submission *datastruct.Submission) error
	Delete(ctx context.Context, id int) error
}

type ValidatorRepository interface {
	ByProblemId(ctx context.Context, problemId int) (*datastruct.Validator, error)
}

type TestCaseRepository interface {
	ByProblemId(ctx context.Context, problemId int) ([]*datastruct.TestCase, error)
}

type ContestRepository interface {
	All(ctx context.Context, query *datastruct.ContestQuery) ([]*datastruct.Contest, error)
	ById(ctx context.Context, id int) (*datastruct.Contest, error)
	Create(ctx context.Context, contest *datastruct.Contest) error
	Update(ctx context.Context, contest *datastruct.Contest) error
	Delete(ctx context.Context, id int) error
}