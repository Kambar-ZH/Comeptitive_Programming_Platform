package test

import (
	"fmt"
	"site/internal/models"
)

const (
	participantSolutionFile = "../../cmd/myapp/participant_solution.go"
	participantSolutionExe  = "../../cmd/myapp/participant_solution.exe"
	mainSolutionFile        = "../../cmd/myapp/main_solution.go"
	mainSolutionExe         = "../../cmd/myapp/main_solution.exe"
	clearFile               = "../../cmd/myapp/clear.txt"
)

var (
	tests = []models.Validator{
		{
			ProblemName:  "0001",
			SolutionFile: "../../test/problems/0001/solution.go",
			TestCases: []models.TestCase{
				{
					InputFile: "../../test/problems/0001/tests/1.txt",
				},
				{
					InputFile: "../../test/problems/0001/tests/2.txt",
				},
			},
		},
		{
			ProblemName:  "0002",
			SolutionFile: "../../test/problems/0002/solution.go",
			TestCases: []models.TestCase{
				{
					InputFile: "../../test/problems/0002/tests/1.txt",
				},
				{
					InputFile: "../../test/problems/0002/tests/2.txt",
				},
			},
		},
	}
)

func RunTestCase(solutionFile, tempFile string, id int, testCase models.TestCase) models.Verdict {
	if err := WriteFile(mainSolutionFile, solutionFile); err != nil {
		fmt.Println("Error on writing to main solution file")
		return models.UNKNOWN_ERROR
	}
	if err := WriteFile(participantSolutionFile, tempFile); err != nil {
		fmt.Println("Error on writing from participant solution file to temp file")
		return models.UNKNOWN_ERROR
	}

	if err := MakeExe(); err != nil {
		fmt.Println("Error on making executable")
		return models.COMPILATION_ERROR
	}

	expected, err := Run(mainSolutionExe, testCase)
	if err != nil {
		fmt.Println("Error on running solution file")
		return models.COMPILATION_ERROR
	}
	actual, err := Run(participantSolutionExe, testCase)
	if err != nil {
		fmt.Println("Error on running participant solution file")
		return models.UNKNOWN_ERROR
	}

	if err = WriteFile(participantSolutionFile, clearFile); err != nil {
		return models.UNKNOWN_ERROR
	}
	if err = WriteFile(mainSolutionFile, clearFile); err != nil {
		return models.UNKNOWN_ERROR
	}

	if expected != actual {
		fmt.Printf("FAIL on test %d.\n\n   Expected: %s\n   Actual: %s\n", id, expected, actual)
		return models.FAILED
	}

	fmt.Printf("OK on test %d.\n\n   Expected: %s\n   Actual: %s\n", id, expected, actual)
	return models.PASSED
}

func TestSolution(tempFile string, problemName string) (models.Verdict, int) {
	for _, test := range tests {
		if test.ProblemName == problemName {
			for id, testCase := range test.TestCases {
				verdict := RunTestCase(test.SolutionFile, tempFile, id+1, testCase)
				if verdict != models.PASSED {
					return verdict, id + 1
				}
			}
		}
	}
	return models.PASSED, -1
}