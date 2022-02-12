package dto

import "time"

type ContestFindAllRequest struct {
	Filter string
	Page   int32
	Limit  int32
	Offset int32
}

type ContestFindByTimeInterval struct {
	StartTime time.Time
	EndTime   time.Time
}

type phase int

const (
	BEFORE phase = iota
	CODING
	PENDING_SYSTEM_TEST
	SYSTEM_TEST
	FINISHED
)

func (p phase) String() string {
	return [...]string{"BEFORE", "CODING", "PENDING_SYSTEM_TEST", "SYSTEM_TEST", "FINISHED"}[p]
}
