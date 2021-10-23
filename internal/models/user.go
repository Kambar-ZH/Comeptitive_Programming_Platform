package models

type User struct {
	Id         int    `json:"id"`
	Handle     string `json:"handle"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Rating     int    `json:"rating"`
	MaxRatring int    `json:"maxRating"`
	Avatar     string `json:"avatar"`
}

type RatingChange struct {
	ContestId int
	Handle    string
	OldRating int
	NewRating int
}

type Contest struct {
	Id              int
	IsFinished      bool
	DurationSeconds int
	Description     string
}

type Problem struct {
	ContestId int
	Id        int
	Name      string
	Rating    int
	Tags      []string
}

type SubmissionResult struct {
	Verdict    Verdict
	FailedTest int
}

type Submission struct {
	Id               int    `json:"id"`
	Date             string `json:"data"`
	AuthorId         int    `json:"authorId"`
	ProblemId        int    `json:"problemId"`
	SubmissionResult `json:"submissionResult"`
}
