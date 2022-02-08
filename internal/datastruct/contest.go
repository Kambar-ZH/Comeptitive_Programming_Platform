package datastruct

import "time"

type Contest struct {
	Id          int32     `json:"id,omitempty" db:"id,omitempty"`
	Name        string    `json:"name,omitempty" db:"name,omitempty"`
	StartDate   time.Time `json:"start_date,omitempty" db:"start_date,omitempty"`
	EndDate     time.Time `json:"end_date,omitempty" db:"end_date,omitempty"`
	Description string    `json:"description,omitempty" db:"description,omitempty"`
	Phase       string    `json:"phase,omitempty" db:"phase,omitempty"`
}
