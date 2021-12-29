package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/middleware"
	"site/internal/store"
)

const (
	usersPerPage = 20
)

type UserService interface {
	All(ctx context.Context, query *datastruct.UserAllRequest) ([]*datastruct.User, error)
	ByEmail(ctx context.Context, email string) (*datastruct.User, error)
	ByHandle(ctx context.Context, handle string) (*datastruct.User, error)
	Create(ctx context.Context, user *datastruct.User) error
	Update(ctx context.Context, user *datastruct.User) error
	Delete(ctx context.Context, handle string) error
}

type UserServiceImpl struct {
	store store.Store
}

func NewUserService(opts ...UserServiceOption) UserService {
	svc := &UserServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (u UserServiceImpl) All(ctx context.Context, req *datastruct.UserAllRequest) ([]*datastruct.User, error) {
	req.Limit = usersPerPage
	req.Offset = (req.Page - 1) * usersPerPage
	users, err := u.store.Users().All(ctx, req)
	return users, err
}

func (u UserServiceImpl) ByEmail(ctx context.Context, email string) (*datastruct.User, error) {
	return u.store.Users().ByEmail(ctx, email)
}

func (u UserServiceImpl) ByHandle(ctx context.Context, handle string) (*datastruct.User, error) {
	return u.store.Users().ByHandle(ctx, handle)
}

func (u UserServiceImpl) Create(ctx context.Context, user *datastruct.User) error {
	if err := middleware.Validate(user); err != nil {
		return err
	}
	if err := middleware.BeforeCreate(user); err != nil {
		return err
	}
	middleware.Sanitize(user)

	return u.store.Users().Create(ctx, user)
}

func (u UserServiceImpl) Update(ctx context.Context, user *datastruct.User) error {
	// TODO: add admin permission
	_, err := u.store.Users().ByHandle(ctx, user.Handle)
	if err != nil {
		return err
	}
	return u.store.Users().Update(ctx, user)
}

func (u UserServiceImpl) Delete(ctx context.Context, handle string) error {
	// TODO: add admin permission
	_, err := u.store.Users().ByHandle(ctx, handle)
	if err != nil {
		return err
	}
	return u.store.Users().Delete(ctx, handle)
}
