package datastruct

type (
	Validator struct {
		ProblemId        int32  `json:"problem_id" db:"problem_id"`
		SolutionFilePath string `json:"solution_file_path" db:"solution_file_path"`
		TestCases        []TestCase `json:"test" db:"solution_file_path"`
	}
)
