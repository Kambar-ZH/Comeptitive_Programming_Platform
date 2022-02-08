package dto

type (
	ProblemResultsFindAllRequest struct {
		ContestId int32
		Filter    string
		Page      int32
		Limit     int32
		Offset    int32
	}

	ProblemResultsGetByProblemIdRequest struct {
		ContestId int32
		UserId    int32
		ProblemId int32
		Filter    string
		Page      int32
		Limit     int32
		Offset    int32
	}
)
