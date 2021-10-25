package compiler

import (
	"fmt"
	"site/internal/models"
	"site/test/inmemory"
)

func PrepareExe(solutionFile, tempFile string) error {
	if err := CopyFile(inmemory.GetInstance().MainSolution, solutionFile); err != nil {
		fmt.Println("Error on writing to main solution")
		return err
	}
	if err := CopyFile(inmemory.GetInstance().ParticipantSolution, tempFile); err != nil {
		fmt.Println("Error on writing to participant solution")
		return err
	}
	if err := BuildExe(); err != nil {
		fmt.Println("Error on making executable")
		return err
	}
	return nil
}

func CleanUp() error {
	if err := CopyFile(inmemory.GetInstance().ParticipantSolution, inmemory.GetInstance().CleanFile); err != nil {
		return err
	}
	if err := CopyFile(inmemory.GetInstance().MainSolution, inmemory.GetInstance().CleanFile); err != nil {
		return err
	}
	return nil
}

func RunTestCase(solutionFile, tempFile string, id int, testCase models.TestCase) models.Verdict {
	expected, err := MustExecuteFile(inmemory.GetInstance().MainSolutionExe, testCase)
	if err != nil {
		fmt.Println("Error on running main solution")
		return models.UNKNOWN_ERROR
	}
	actual, err := MustExecuteFile(inmemory.GetInstance().ParticipantSolutionExe, testCase)
	if err != nil {
		fmt.Println("Error on running participant solution")
		return models.COMPILATION_ERROR
	}

	if expected != actual {
		fmt.Printf("FAIL on test %d:\nExpected: %s\nActual: %s\n", id, expected, actual)
		return models.FAILED
	}

	fmt.Printf("OK on test %d:\nExpected: %s\nActual: %s\n", id, expected, actual)
	return models.PASSED
}

func TestSolution(tempFile string, problemId int) models.SubmissionResult {
	test, err := inmemory.GetInstance().GetTestByID(problemId)
	if err != nil {
		return models.SubmissionResult{
			Verdict: models.UNKNOWN_ERROR,
		}
	}
	if err := PrepareExe(test.SolutionFile, tempFile); err != nil {
		return models.SubmissionResult{
			Verdict: models.COMPILATION_ERROR,
		}
	}
	defer CleanUp()
	for id, testCase := range test.TestCases {
		verdict := RunTestCase(test.SolutionFile, tempFile, id+1, testCase)
		if verdict != models.PASSED {
			return models.SubmissionResult{
				Verdict:    verdict,
				FailedTest: testCase.Id + 1,
			}
		}
	}
	return models.SubmissionResult{
		Verdict: models.PASSED,
	}
}