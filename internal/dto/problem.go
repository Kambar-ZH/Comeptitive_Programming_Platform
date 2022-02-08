package dto

import "site/internal/datastruct"

type (
	ProblemFindAllRequest struct {
		Page      int32
		Limit     int32
		Offset    int32
		ContestId int32
	}

	ProblemsetRequest struct {
		Page          int32
		Limit         int32
		Offset        int32
		MinDifficulty int32
		MaxDifficulty int32
		FilterTag     string
	}

	ProblemGetByIdRequest struct {
		ContestId int32
		ProblemId int32
	}

	ProblemUpdateRequest struct {
		ContestId int32
		Problem   *datastruct.Problem
	}

	ProblemDeleteRequest struct {
		ContestId int32
		ProblemId int32
	}

	ProblemCreateRequest struct {
		ContestId int32
		Problem   *datastruct.Problem
	}
)
