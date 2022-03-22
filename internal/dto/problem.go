package dto

import (
	"site/internal/consts"
	"site/internal/datastruct"
)

type (
	ProblemFindAllRequest struct {
		LanguageCode consts.LanguageCode
		ContestId    int32
	}

	ProblemsetRequest struct {
		LanguageCode  consts.LanguageCode
		MinDifficulty int32
		MaxDifficulty int32
		Pagination
	}

	ProblemGetByIdRequest struct {
		ProblemId    int32
		ContestId    int32
		LanguageCode consts.LanguageCode
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
