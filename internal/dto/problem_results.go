package dto

type (
	ProblemResultsFindAllRequest struct {
		ContestId int32
	}

	ProblemResultGetByProblemIdRequest struct {
		ContestId int32
		UserId    int32
		ProblemId int32
	}
)
