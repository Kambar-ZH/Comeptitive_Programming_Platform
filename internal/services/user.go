package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/store"
)

const (
	usersPerPage = 20
)

type UserService interface {
	FindAll(ctx context.Context, req *dto.UserFindAllRequest) ([]*datastruct.User, error)
	GetByEmail(ctx context.Context, email string) (*datastruct.User, error)
	GetByHandle(ctx context.Context, handle string) (*datastruct.User, error)
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

func (u UserServiceImpl) FindAll(ctx context.Context, req *dto.UserFindAllRequest) ([]*datastruct.User, error) {
	req.Limit = usersPerPage
	req.Offset = (req.Page - 1) * usersPerPage
	users, err := u.store.Users().FindAll(ctx, req)
	return users, err
}

func (u UserServiceImpl) GetByEmail(ctx context.Context, email string) (*datastruct.User, error) {
	return u.store.Users().GetByEmail(ctx, email)
}

func (u UserServiceImpl) GetByHandle(ctx context.Context, handle string) (*datastruct.User, error) {
	return u.store.Users().GetByHandle(ctx, handle)
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
	_, err := u.store.Users().GetByHandle(ctx, user.Handle)
	if err != nil {
		return err
	}
	return u.store.Users().Update(ctx, user)
}

func (u UserServiceImpl) Delete(ctx context.Context, handle string) error {
	// TODO: add admin permission
	_, err := u.store.Users().GetByHandle(ctx, handle)
	if err != nil {
		return err
	}
	return u.store.Users().Delete(ctx, handle)
}
