package datastruct

type (
	TestCase struct {
		Id            int32  `json:"id,omitempty" db:"id,omitempty"`
		ProblemId     int32  `json:"problem_id,omitempty" db:"problem_id,omitempty"`
		AsSample      bool   `json:"as_sample,omitempty" db:"as_sample,omitempty"`
		InputFilePath string `json:"input_file_path,omitempty" db:"input_file_path,omitempty"`
	}
)
