package consts

type Verdict int

const (
	PASSED Verdict = iota
	FAILED
	COMPILATION_ERROR
	TIME_LIMIT_ERROR
	UNKNOWN_ERROR
)

func (v Verdict) String() string {
	return [...]string{"Passed", "Failed Test", "Compilation Error", "Time Limit Error", "Unknown Error"}[v]
}
