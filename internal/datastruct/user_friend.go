package datastruct

type UserFriend struct {
	UserId   int32 `json:"user_id" db:"user_id"`
	FriendId int32 `json:"friends_id" db:"friends_id"`
}
