package inmemory

import (
	"context"
	"fmt"
	"site/internal/grpc/api"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepo struct {
	data map[int32]*api.User
	api.UnimplementedUserRepositoryServer
	mu *sync.RWMutex
}

func (db *UserRepo) All(ctx context.Context, empty *api.Empty) (*api.UserList, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := []*api.User{}
	for _, user := range db.data {
		res = append(res, user)
	}
	ans := api.UserList{Users: res}

	return &ans, nil
}

func (db *UserRepo) ById(ctx context.Context, user *api.UserRequestId) (*api.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if user, ok := db.data[user.Id]; ok {
		return user, nil
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with id %d does not exist", user.Id))
}

func (db *UserRepo) Create(ctx context.Context, user *api.User) (*api.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[user.Id] = user
	return user, nil
}

func (db *UserRepo) Update(ctx context.Context, user *api.User) (*api.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[user.Id] = user

	return user, nil
}

func (db *UserRepo) Delete(ctx context.Context, user *api.UserRequestId) (*api.Empty, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.data[user.Id]; !ok {
		return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
	}
	delete(db.data, user.Id)
	return &api.Empty{}, nil
}
