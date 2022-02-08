package dto

import "site/internal/datastruct"

type (
	SubmissionFindAllRequest struct {
		FilterUserHandle string
		Page             int32
		Limit            int32
		Offset           int32
		ContestId        int32
	}

	SubmissionCreateRequest struct {
		Submission *datastruct.Submission
		ContestId  int32
	}

	SubmissionUpdateRequest struct {
		Submission *datastruct.Submission
		ContestId  int32
	}

	SubmissionDeleteRequest struct {
		SubmissionId int32
		ContestId    int32
	}

	SubmissionGetByIdRequest struct {
		SubmissionId int32
		ContestId    int32
	}
)
