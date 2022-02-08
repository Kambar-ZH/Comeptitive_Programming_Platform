package datastruct

type Participant struct {
	UserId          int32  `json:"user_id,omitempty" db:"user_id,omitempty"`
	ContestId       int32  `json:"contest_id,omitempty" db:"contest_id,omitempty"`
	ParticipantType string `json:"participant_type,omitempty" db:"participant_type,omitempty"`
	Room            int32  `json:"room,omitempty" db:"room,omitempty"`
	ProblemResults  []ProblemResults `json:"problem_results,omitempty" db:"problem_results,omitempty"`
}
