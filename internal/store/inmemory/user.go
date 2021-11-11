package inmemory

import (
	"context"
	"fmt"
	"site/internal/datastruct"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepo struct {
	data map[string]*datastruct.User
	mu   *sync.RWMutex
}

// db *UserRepo

func (db *UserRepo) All(ctx context.Context, offset, limit int) ([]*datastruct.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := []*datastruct.User{}
	for _, user := range db.data {
		res = append(res, user)
	}

	return res, nil
}

func (db *UserRepo) ByEmail(ctx context.Context, email string) (*datastruct.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, user := range db.data {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with email %s does not exist", email))
}

func (db *UserRepo) ByHandle(ctx context.Context, handle string) (*datastruct.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if user, ok := db.data[handle]; ok {
		return user, nil
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with handle %s does not exist", handle))
}

func (db *UserRepo) Create(ctx context.Context, user *datastruct.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[user.Handle] = user
	return nil
}

func (db *UserRepo) Update(ctx context.Context, user *datastruct.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[user.Handle] = user
	return nil
}

func (db *UserRepo) Delete(ctx context.Context, handle string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, handle)
	return nil
}