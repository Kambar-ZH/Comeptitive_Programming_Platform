package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"site/internal/datastruct"
	"site/internal/store"
)

func (db *DB) UserFriends() store.UserFriendRepository {
	if db.userFriends == nil {
		db.userFriends = NewUserFriendRepository(db.conn)
	}
	return db.userFriends
}

type UserFriendRepository struct {
	conn *sqlx.DB
}

func NewUserFriendRepository(conn *sqlx.DB) store.UserFriendRepository {
	return &UserFriendRepository{conn: conn}
}

func (u UserFriendRepository) Create(ctx context.Context, userFriend *datastruct.UserFriend) error {
	if _, err := u.conn.Exec(
		`INSERT INTO user_friends(user_id, friend_id) 
    		VALUES ($1, $2)`,
		userFriend.UserId, userFriend.FriendId); err != nil {
		return err
	}
	return nil
}

func (u UserFriendRepository) Delete(ctx context.Context, userFriend *datastruct.UserFriend) error {
	if _, err := u.conn.Exec(
		`DELETE FROM user_friends
			WHERE user_id = $1
			AND friend_id = $2`,
		userFriend.UserId, userFriend.FriendId); err != nil {
		return err
	}
	return nil
}
