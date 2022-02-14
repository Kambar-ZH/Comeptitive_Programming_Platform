package dto

type UserFindAllRequest struct {
	Filter string
	Page   int32
	Limit  int32
	Offset int32
}

type UserFindFriendsRequest struct {
	UserId int32
	Page   int32
	Limit  int32
	Offset int32
}
