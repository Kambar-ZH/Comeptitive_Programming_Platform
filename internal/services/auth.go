package services

import (
	"context"
	"errors"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/store"
)

var (
	errInvalidEmailOrPassword = errors.New("invalid email or password")
)

type AuthService interface {
	ByEmail(ctx context.Context, req *dto.Cridentials) (*datastruct.User, error)
	ByHandle(ctx context.Context, handle string) (*datastruct.User, error)
}

type AuthServiceImpl struct {
	store store.Store
}

func NewAuthService(opts ...AuthServiceOption) AuthService {
	svc := &AuthServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (a AuthServiceImpl) ByEmail(ctx context.Context, req *dto.Cridentials) (*datastruct.User, error) {
	user, err := a.store.Users().GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if !middleware.ComparePassword(user, req.Password) {
		return nil, errInvalidEmailOrPassword
	}
	return user, nil
}

func (a AuthServiceImpl) ByHandle(ctx context.Context, handle string) (*datastruct.User, error) {
	return a.store.Users().GetByHandle(ctx, handle)
}