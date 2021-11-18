package datastruct

import "time"

type (
	Contest struct {
		Id           int32     `json:"id,omitempty" db:"id,omitempty"`
		Name         string    `json:"name,omitempty" db:"name,omitempty"`
		StartDate    time.Time `json:"start_date,omitempty" db:"start_date,omitempty"`
		Description  string    `json:"description,omitempty" db:"description,omitempty"`
		AuthorHandle string    `json:"author_handle,omitempty" db:"author_handle,omitempty"`
		Problems     []Problem `json:"problems,omitempty" db:"problems,omitempty"`
	}
)
