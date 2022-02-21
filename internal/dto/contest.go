package dto

import (
	"site/internal/consts"
	"time"
)

type ContestFindAllRequest struct {
	Filter       string
	Page         int32
	Limit        int32
	Offset       int32
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
