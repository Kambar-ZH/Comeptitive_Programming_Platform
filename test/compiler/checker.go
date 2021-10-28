package compiler

import (
	"fmt"
	"site/grpc/api"
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

func RunTestCase(solutionFile, tempFile string, id int, testCase models.TestCase) api.Verdict {
	expected, err := MustExecuteFile(inmemory.GetInstance().MainSolutionExe, testCase)
	if err != nil {
		fmt.Println("Error on running main solution")
		return api.Verdict_UNKNOWN_ERROR
	}
	actual, err := MustExecuteFile(inmemory.GetInstance().ParticipantSolutionExe, testCase)
	if err != nil {
		fmt.Println("Error on running participant solution")
		return api.Verdict_COMPILATION_ERROR
	}

	if expected != actual {
		fmt.Printf("[%d] incorrect result on test::\nExpected: %s\nActual: %s\n", id, expected, actual)
		return api.Verdict_FAILED
	}

	fmt.Printf("[%d] correct result on test:\nExpected: %s\nActual: %s\n", id, expected, actual)
	return api.Verdict_PASSED
}

func TestSolution(tempFile string, problemId int) *api.SubmissionResult {
	test, err := inmemory.GetInstance().GetTestByID(problemId)
	if err != nil {
		return &api.SubmissionResult{
			Verdict: api.Verdict_UNKNOWN_ERROR,
		}
	}
	if err := PrepareExe(test.SolutionFile, tempFile); err != nil {
		return &api.SubmissionResult{
			Verdict: api.Verdict_COMPILATION_ERROR,
		}
	}
	defer CleanUp()
	for id, testCase := range test.TestCases {
		verdict := RunTestCase(test.SolutionFile, tempFile, id+1, testCase)
		if verdict != api.Verdict_PASSED {
			return &api.SubmissionResult{
				Verdict:    verdict,
				FailedTest: int32(testCase.Id + 1),
			}
		}
	}
	return &api.SubmissionResult{
		Verdict: api.Verdict_PASSED,
	}
}