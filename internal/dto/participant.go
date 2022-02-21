package dto

type ParticipantFindAllRequest struct {
	ContestId int
	Filter    string
	Page      int32
	Limit     int32
	Offset    int32
}

type ParticipantFindFriendsRequest struct {
	ContestId int
	UserId    int
	Filter    string
	Page      int32
	Limit     int32
	Offset    int32
}

type ParticipantGetByUserIdRequest struct {
	ContestId int
	UserId    int
}

type ParticipantRegisterRequest struct {
	ContestId       int
	ParticipantType string
}
