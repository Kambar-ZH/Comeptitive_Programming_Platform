package inmemory

// import (
// 	"site/internal/datastruct"
// 	"site/internal/store"
// 	"sync"
// )

// type DB struct {
// 	usersRepo       store.UserRepository
// 	submissionsRepo store.SubmissionRepository

// 	mu *sync.RWMutex
// }

// func NewDB() store.Store {
// 	return &DB{
// 		mu: new(sync.RWMutex),
// 	}
// }

// // redudent for inmemory storage
// func (db *DB) Connect(url string) error {
// 	return nil
// }

// func (db *DB) Close() error {
// 	return nil
// }

// func (db *DB) Users() store.UserRepository {
// 	if db.usersRepo == nil {
// 		db.usersRepo = &UserRepo{
// 			data: make(map[string]*datastruct.User),
// 			mu:   &sync.RWMutex{},
// 		}
// 	}

// 	return db.usersRepo
// }

// func (db *DB) Submissions() store.SubmissionRepository {
// 	if db.submissionsRepo == nil {
// 		db.submissionsRepo = &SubmissionRepo{
// 			data: make(map[int32]*datastruct.Submission),
// 			mu:   &sync.RWMutex{},
// 		}
// 	}

// 	return db.submissionsRepo
// }
