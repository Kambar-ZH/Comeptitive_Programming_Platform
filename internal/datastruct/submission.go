package datastruct

type Submission struct {
	Id           int32  `json:"id,omitempty"`
	ContestId    int32  `json:"contestId,omitempty"`
	Date         string `json:"date,omitempty"`
	AuthorHandle string `json:"authorHandle,omitempty"`
	ProblemId    int32  `json:"problemId,omitempty"`
	Verdict      string  `json:"verdict,omitempty"`
	FailedTest   int32  `json:"failedTest,omitempty"`
}