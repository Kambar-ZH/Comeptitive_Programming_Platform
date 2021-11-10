package inmemory

import (
	"context"
	"fmt"
	"site/internal/grpc/api"
	"site/internal/services"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepo struct {
	data map[string]*api.User
	api.UnimplementedUserRepositoryServer
	mu *sync.RWMutex
}

func (db *UserRepo) All(ctx context.Context, empty *api.AllUserRequest) (*api.UserList, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := []*api.User{}
	for _, user := range db.data {
		res = append(res, user)
	}
	ans := api.UserList{Users: res}

	return &ans, nil
}

func (db *UserRepo) ByHandle(ctx context.Context, req *api.UserByHandleRequest) (*api.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if user, ok := db.data[req.Handle]; ok {
		return user, nil
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with handle %s does not exist", req.Handle))
}

func (db *UserRepo) ByEmail(ctx context.Context, req *api.UserByEmailRequest) (*api.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, user := range db.data {
		if user.Email == req.Email {
			return user, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with email %s does not exist", req.Email))
}

func (db *UserRepo) Create(ctx context.Context, user *api.User) (*api.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := services.Validate(user); err != nil {
		return nil, err
	}

	if err := services.BeforeCreate(user); err != nil {
		return nil, err
	}

	db.data[user.Handle] = user
	return user, nil
}

func (db *UserRepo) Update(ctx context.Context, user *api.User) (*api.User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[user.Handle] = user

	return user, nil
}

func (db *UserRepo) Delete(ctx context.Context, req *api.DeleteUserRequest) (*api.Empty, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.data[req.Handle]; !ok {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("user with handle %s does not exist", req.Handle))
	}
	delete(db.data, req.Handle)
	return &api.Empty{}, nil
}
