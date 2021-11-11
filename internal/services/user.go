package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"
)

const (
	usersPerPage = 20
)

type UserService interface {
	All(ctx context.Context, page int) ([]*datastruct.User, error)
	ByEmail(ctx context.Context, email string) (*datastruct.User, error)
	ByHandle(ctx context.Context, handle string) (*datastruct.User, error)
	Create(ctx context.Context, user *datastruct.User) error
	Update(ctx context.Context, user *datastruct.User) error
	Delete(ctx context.Context, handle string) error
}

type UserServiceImpl struct {
	repo store.UserRepository
}

func NewUserService(opts ...UserServiceOption) UserService {
	svc := &UserServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (u UserServiceImpl) All(ctx context.Context, page int) ([]*datastruct.User, error) {
	return u.repo.All(ctx, (page-1)*usersPerPage, usersPerPage)
}

func (u UserServiceImpl) ByEmail(ctx context.Context, email string) (*datastruct.User, error) {
	return u.repo.ByEmail(ctx, email)
}

func (u UserServiceImpl) ByHandle(ctx context.Context, handle string) (*datastruct.User, error) {
	return u.repo.ByHandle(ctx, handle)
}

func (u UserServiceImpl) Create(ctx context.Context, user *datastruct.User) error {
	// TODO: VALIDATION
	// if err := services.Validate(user); err != nil {
	// 	return nil, err
	// }

	// if err := services.BeforeCreate(user); err != nil {
	// 	return nil, err
	// }
	
	// services.Sanitize(user)
	return u.repo.Create(ctx, user)
}

func (u UserServiceImpl) Update(ctx context.Context, user *datastruct.User) error {
	// TODO: add admin permission
	_, err := u.repo.ByHandle(ctx, user.Handle)
	if err != nil {
		return err
	}
	return u.repo.Update(ctx, user)
}

func (u UserServiceImpl) Delete(ctx context.Context, handle string) error {
	// TODO: add admin permission
	_, err := u.repo.ByHandle(ctx, handle)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, handle)
}
