package inmemory

import (
	"site/grpc/api"
	"site/internal/store"
	"sync"
)

type DB struct {
	usersRepo       api.UserRepositoryServer
	submissionsRepo api.SubmissionRepositoryServer

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Users() store.UserRepository {
	if db.usersRepo == nil {
		db.usersRepo = &UserRepo{
			data: make(map[int32]*api.User),
			mu:   &sync.RWMutex{},
		}
	}

	return db.usersRepo
}

func (db *DB) Submissions() store.SubmissionRepository {
	if db.submissionsRepo == nil {
		db.submissionsRepo = &SubmissionRepo{
			data: make(map[int32]*api.Submission),
			mu:   &sync.RWMutex{},
		}
	}

	return db.submissionsRepo
}
