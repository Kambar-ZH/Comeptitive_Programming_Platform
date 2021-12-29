package datastruct

type (
	Problem struct {
		Id         int32  `json:"id,omitempty" db:"id,omitempty"`
		ContestId  int32  `json:"contest_id,omitempty" db:"contest_id,omitempty"`
		Index      string `json:"index,omitempty" db:"index,omitempty"`
		Name       string `json:"name,omitempty" db:"name,omitempty"`
		Statement  string `json:"statement,omitempty" db:"statement,omitempty"`
		Input      string `json:"input,omitempty" db:"input,omitempty"`
		Output     string `json:"output,omitempty" db:"output,omitempty"`
		Difficulty int32  `json:"difficulty,omitempty" db:"difficulty,omitempty"`
		Tags       []Tag
	}

	ProblemAllRequest struct {
		Page      int32
		Limit     int32
		Offset    int32
		ContestId int32
	}

	ProblemsetRequest struct {
		Page          int32
		Limit         int32
		Offset        int32
		MinDifficulty int32
		MaxDifficulty int32
		FilterTag     string
	}

	ProblemByIdRequest struct {
		ContestId int32
		ProblemId int32
	}

	ProblemUpdateRequest struct {
		ContestId int32
		Problem   *Problem
	}

	ProblemDeleteRequest struct {
		ContestId int32
		ProblemId int32
	}

	ProblemCreateRequest struct {
		ContestId int32
		Problem   *Problem
	}
)
