package consts

type Phase int

const (
	BEFORE Phase = iota
	CODING
	PENDING_SYSTEM_TEST
	SYSTEM_TEST
	FINISHED
)

func (p Phase) String() string {
	return [...]string{"BEFORE", "CODING", "PENDING_SYSTEM_TEST", "SYSTEM_TEST", "FINISHED"}[p]
}
