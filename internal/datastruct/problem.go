package datastruct

type Problem struct {
	Id         int32  `json:"id,omitempty" db:"id,omitempty"`
	ContestId  int32  `json:"contest_id,omitempty" db:"contest_id,omitempty"`
	Index      string `json:"index,omitempty" db:"index,omitempty"`
	Name       string `json:"name,omitempty" db:"name,omitempty"`
	Statement  string `json:"statement,omitempty" db:"statement,omitempty"`
	Input      string `json:"input,omitempty" db:"input,omitempty"`
	Output     string `json:"output,omitempty" db:"output,omitempty"`
	Difficulty int32  `json:"difficulty,omitempty" db:"difficulty,omitempty"`
	Points     int32  `json:"points,omitempty" db:"points,omitempty"`
	Tags       []*Tag
}
