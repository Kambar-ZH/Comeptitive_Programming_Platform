package datastruct

type Participant struct {
	UserId          int32  `json:"user_id,omitempty" db:"user_id,omitempty"`
	ContestId       int32  `json:"contest_id,omitempty" db:"contest_id,omitempty"`
	ParticipantType string `json:"participant_type,omitempty" db:"participant_type,omitempty"`
	Room            int32  `json:"room,omitempty" db:"room,omitempty"`
	Handle          string `json:"handle,omitempty" db:"handle,omitempty"`
	Rating          int32  `json:"rating,omitempty" db:"rating,omitempty"`
	Points          int32  `json:"points,omitempty" db:"points,omitempty"`
}
