package services

import (
	"context"
	"site/internal/grpc/api"
	"site/internal/store"
)

const (
	usersLimitPerPage = 20
)

type UserService interface {
	All(ctx context.Context, req *api.AllUserRequest) (*api.UserList, error)
	ByEmail(ctx context.Context, req *api.UserByEmailRequest) (*api.User, error)
	ByHandle(ctx context.Context, req *api.UserByHandleRequest) (*api.User, error) 
	Create(ctx context.Context, user *api.User) (*api.User, error) 
	Update(ctx context.Context, user *api.User) (*api.User, error)
	Delete(ctx context.Context, user *api.DeleteUserRequest) (*api.Empty, error)
}

type UserServiceImpl struct {
	repo store.UserRepository
}

func NewUserService(opts ...UserServiceOption) UserService {
	svc := &UserServiceImpl{}
	for _, v := range(opts) {
		v(svc)
	}
	return svc
}

func (u UserServiceImpl) All(ctx context.Context, req *api.AllUserRequest) (*api.UserList, error) {
	req.Limit = usersLimitPerPage

	return u.repo.All(ctx, req)
}

func (u UserServiceImpl) ByEmail(ctx context.Context, req *api.UserByEmailRequest) (*api.User, error) {
	return u.repo.ByEmail(ctx, req)
}

func (u UserServiceImpl) ByHandle(ctx context.Context, req *api.UserByHandleRequest) (*api.User, error) {
	return u.repo.ByHandle(ctx, req)
}

func (u UserServiceImpl) Create(ctx context.Context, user *api.User) (*api.User, error) {
	return u.repo.Create(ctx, user)
}

func (u UserServiceImpl) Update(ctx context.Context, user *api.User) (*api.User, error) {
	// TODO: add admin permission
	req := &api.UserByHandleRequest{
		Handle: user.Handle,
	}
	_, err := u.repo.ByHandle(ctx, req)
	if err != nil {
		return nil, err
	}
	return u.repo.Update(ctx, user)
}

func (u UserServiceImpl) Delete(ctx context.Context, user *api.DeleteUserRequest) (*api.Empty, error) {
	// TODO: add admin permission
	req := &api.UserByHandleRequest{
		Handle: user.Handle,
	}
	_, err := u.repo.ByHandle(ctx, req)
	if err != nil {
		return nil, err
	}
	return u.repo.Delete(ctx, user)
}