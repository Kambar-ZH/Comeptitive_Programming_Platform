package datastruct

type Validator struct {
	ProblemId              int32       `json:"problem_id" db:"problem_id"`
	AuthorSolutionFilePath string      `json:"solution_file" db:"solution_file"`
	TestCases              []*TestCase `json:"test_cases" db:"test_cases"`
}
