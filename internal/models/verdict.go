package models

type Verdict int64

const (
	PASSED Verdict = iota
	FAILED
	COMPILATION_ERROR
	TIME_LIMIT_ERROR
	UNKNOWN_ERROR
)

func (v Verdict) String() string {
	switch v {
	case PASSED:
		return "Passed"
	case FAILED:
		return "Failed"
	case COMPILATION_ERROR:
		return "Compilation Error"
	case TIME_LIMIT_ERROR:
		return "Time limit error"
	}
	return "Unknown error"
}