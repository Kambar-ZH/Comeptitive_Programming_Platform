package dto

import (
	"site/internal/consts"
	"time"
)

type ContestFindAllRequest struct {
	Pagination
	LanguageCode consts.LanguageCode
}

type ContestFindByTimeIntervalRequest struct {
	StartTime    time.Time
	EndTime      time.Time
	LanguageCode consts.LanguageCode
}

type ContestGetByIdRequest struct {
	ContestId    int32
	LanguageCode consts.LanguageCode
}
