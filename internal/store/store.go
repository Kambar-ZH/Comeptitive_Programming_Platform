package store

import (
	"context"
	"site/internal/models"
)

type Store interface {
	Users() UserRepository
	Submissions() SubmissionRepository
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context) ([]*models.User, error)
	ById(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

type SubmissionRepository interface {
	Create(ctx context.Context, submission *models.Submission) error
	All(ctx context.Context) ([]*models.Submission, error)
	ById(ctx context.Context, id int) (*models.Submission, error)
	Update(ctx context.Context, submission *models.Submission) error
	Delete(ctx context.Context, id int) error
}