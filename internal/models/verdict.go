package models

type Verdict string

const (
	PASSED Verdict = "Passed"
	FAILED Verdict = "Failed Test"
	COMPILATION_ERROR = "Compilation Error"
	TIME_LIMIT_ERROR = "Time Limit Error"
	UNKNOWN_ERROR = "Unknown Error"
)