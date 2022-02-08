package datastruct

import "time"

type Submission struct {
	Id               int32     `json:"id,omitempty" db:"id,omitempty"`
	UserId           int32     `json:"user_id,omitempty" db:"user_id,omitempty"`
	ContestId        int32     `json:"contest_id,omitempty" db:"contest_id,omitempty"`
	ProblemId        int32     `json:"problem_id,omitempty" db:"problem_id,omitempty"`
	SubmissionTime   time.Time `json:"submission_time,omitempty" db:"submission_time,omitempty"`
	Verdict          string    `json:"verdict,omitempty" db:"verdict,omitempty"`
	FailedTest       int32     `json:"failed_test,omitempty" db:"failed_test,omitempty"`
	SolutionFilePath string    `json:"solution_file_path,omitempty" db:"solution_file_path,omitempty"`
}
