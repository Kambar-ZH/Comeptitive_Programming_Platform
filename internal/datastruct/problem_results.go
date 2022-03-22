package datastruct

import "time"

type ProblemResult struct {
	UserId                       int32     `json:"user_id,omitempty" db:"user_id,omitempty"`
	ContestId                    int32     `json:"contest_id,omitempty" db:"contest_id,omitempty"`
	ProblemId                    int32     `json:"problem_id,omitempty" db:"problem_id,omitempty"`
	Points                       int32     `json:"points,omitempty" db:"points,omitempty"`
	Penalty                      int32     `json:"penalty,omitempty" db:"penalty,omitempty"`
	LastSuccessfulSubmissionTime time.Time `json:"last_successful_submission_time,omitempty" db:"last_successful_submission_time,omitempty"`
}
