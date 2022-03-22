package dto

import "site/internal/datastruct"

type (
	SubmissionFindAllRequest struct {
		ContestId int32
		Pagination
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
