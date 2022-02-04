package datastruct

import "time"

type (
	Submission struct {
		Id             int32     `json:"id,omitempty" db:"id,omitempty"`
		ContestId      int32     `json:"contest_id,omitempty" db:"contest_id,omitempty"`
		SubmissionDate time.Time `json:"submission_date,omitempty" db:"submission_date,omitempty"`
		UserId         int32     `json:"user_id,omitempty" db:"user_id,omitempty"`
		ProblemId      int32     `json:"problem_id,omitempty" db:"problem_id,omitempty"`
		Verdict        string    `json:"verdict,omitempty" db:"verdict,omitempty"`
		FailedTest     int32     `json:"failed_test,omitempty" db:"failed_test,omitempty"`
	}

	SubmissionAllRequest struct {
		FilterUserHandle string
		Page             int32
		Limit            int32
		Offset           int32
		ContestId        int32
	}

	SubmissionCreateRequest struct {
		Submission *Submission
		ContestId  int32
	}

	SubmissionUpdateRequest struct {
		Submission *Submission
		ContestId  int32
	}

	SubmissionDeleteRequest struct {
		SubmissionId int32
		ContestId    int32
	}

	SubmissionByIdRequest struct {
		SubmissionId int32
		ContestId    int32
	}
)
