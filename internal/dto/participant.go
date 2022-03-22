package dto

type ParticipantFindAllRequest struct {
	ContestId int32
	Pagination
}

type ParticipantFindFriendsRequest struct {
	ContestId int32
	UserId    int32
	Pagination
}

type ParticipantGetByIdRequest struct {
	ContestId int32
	UserId    int32
}

type ParticipantRegisterRequest struct {
	ContestId       int32
	ParticipantType string
}

type GetStandingsRequest struct {
	ContestId int32
	UserId    int32
	Pagination
}
