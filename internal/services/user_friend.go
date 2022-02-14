package services

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/store"
)

type UserFriendsService interface {
	Create(ctx context.Context, req *dto.UserFriendCreateRequest) error
	Delete(ctx context.Context, req *dto.UserFriendDeleteRequest) error
}

type UserFriendServiceImpl struct {
	store store.Store
}

func NewUserFriendService(opts ...UserFriendServiceOption) UserFriendsService {
	svc := &UserFriendServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (u UserFriendServiceImpl) Create(ctx context.Context, req *dto.UserFriendCreateRequest) error {
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}
	friend, err := u.store.Users().GetByHandle(ctx, req.FriendHandle)
	if err != nil {
		return err
	}
	return u.store.UserFriends().Create(ctx, &datastruct.UserFriend{
		UserId:   user.Id,
		FriendId: friend.Id,
	})
}

func (u UserFriendServiceImpl) Delete(ctx context.Context, req *dto.UserFriendDeleteRequest) error {
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}
	friend, err := u.store.Users().GetByHandle(ctx, req.FriendHandle)
	if err != nil {
		return err
	}
	return u.store.UserFriends().Delete(ctx, &datastruct.UserFriend{
		UserId:   user.Id,
		FriendId: friend.Id,
	})
}
