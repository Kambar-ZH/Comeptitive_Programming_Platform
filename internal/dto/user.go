package dto

type UserFindAllRequest struct {
	Pagination
}

type UserFindFriendsRequest struct {
	UserId int32
	Pagination
}
