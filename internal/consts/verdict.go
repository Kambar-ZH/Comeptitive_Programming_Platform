package consts

type Verdict int

const (
	PRETESTS_PASSED Verdict = iota
	SYSTEM_TESTS_PASSED
	FAILED
	COMPILATION_ERROR
	TIME_LIMIT_ERROR
	UNKNOWN_ERROR
)

func (v Verdict) String() string {
	return [...]string{"Pretests Passed", "System Tests Passed", "Failed Test", "Compilation Error", "Time Limit Error", "Unknown Error"}[v]
}
