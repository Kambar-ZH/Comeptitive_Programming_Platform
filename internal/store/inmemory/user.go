package inmemory

import (
	"context"
	"fmt"
	"site/internal/models"
	"sync"
)

type UsersRepo struct {
	data map[int]*models.User
	mu   *sync.RWMutex
}

func (db *UsersRepo) Create(ctx context.Context, user *models.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[user.Id] = user
	return nil
}

func (db *UsersRepo) All(ctx context.Context) ([]*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	users := make([]*models.User, 0, len(db.data))
	for _, user := range db.data {
		users = append(users, user)
	}

	return users, nil
}

func (db *UsersRepo) ById(ctx context.Context, id int) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no user with id %d", id)
	}
	return user, nil
}

func (db *UsersRepo) Update(ctx context.Context, user *models.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[user.Id] = user
	return nil
}

func (db *UsersRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
