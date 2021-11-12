package datastruct

import "time"

type Submission struct {
	Id             int32     `json:"id,omitempty" db:"id,omitempty"`
	ContestId      int32     `json:"contest_id,omitempty" db:"contest_id,omitempty"`
	SubmissionDate time.Time `json:"submission_date,omitempty" db:"submission_date,omitempty"`
	AuthorHandle   string    `json:"author_handle,omitempty" db:"author_handle,omitempty"`
	ProblemId      int32     `json:"problem_id,omitempty" db:"problem_id,omitempty"`
	Verdict        string    `json:"verdict,omitempty" db:"verdict,omitempty"`
	FailedTest     int32     `json:"failed_test,omitempty" db:"failed_test,omitempty"`
}
