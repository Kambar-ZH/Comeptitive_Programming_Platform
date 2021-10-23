package models

type TestCase struct {
	Id        int
	InputFile string
}


type Validator struct {
	ProblemId  int
	SolutionFile string
	TestCases    []TestCase
}