package inmemory

import (
	"site/internal/models"
	"site/internal/store"
	"sync"
)

type DB struct {
	usersRepo       store.UserRepository
	submissionsRepo store.SubmissionRepository

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Users() store.UserRepository {
	if db.usersRepo == nil {
		db.usersRepo = &UsersRepo{
			data: make(map[int]*models.User),
			mu:   new(sync.RWMutex),
		}
	}

	return db.usersRepo
}

func (db *DB) Submissions() store.SubmissionRepository {
	if db.submissionsRepo == nil {
		db.submissionsRepo = &SubmissionsRepo{
			data: make(map[int]*models.Submission),
			mu:   new(sync.RWMutex),
		}
	}

	return db.submissionsRepo
}