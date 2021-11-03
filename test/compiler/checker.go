package compiler

import (
	"log"
	"site/internal/grpc/api"
	"site/internal/models"
	"site/test/inmemory"
)

func PrepareExe(solutionFile, tempFile string) error {
	if err := CopyFile(inmemory.GetInstance().GetMainSolution(), solutionFile); err != nil {
		log.Println("error on writing to main solution")
		return err
	}
	if err := CopyFile(inmemory.GetInstance().GetParticipantSolution(), tempFile); err != nil {
		log.Println("error on writing to participant solution")
		return err
	}
	if err := BuildExe(); err != nil {
		log.Println("error on making executable")
		return err
	}
	return nil
}

func CleanUp() error {
	if err := CopyFile(inmemory.GetInstance().GetParticipantSolution(), inmemory.GetInstance().GetCleanFile()); err != nil {
		return err
	}
	if err := CopyFile(inmemory.GetInstance().GetMainSolution(), inmemory.GetInstance().GetCleanFile()); err != nil {
		return err
	}
	return nil
}

func RunTestCase(solutionFile, tempFile string, id int, testCase models.TestCase) api.VERDICT {
	expected, err := MustExecuteFile(inmemory.GetInstance().GetMainSolutionExe(), testCase)
	if err != nil {
		log.Println("error on running main solution")
		return api.VERDICT_UNKNOWN_ERROR
	}
	actual, err := MustExecuteFile(inmemory.GetInstance().GetParticipantSolutionExe(), testCase)
	if err != nil {
		log.Println("error on running participant solution")
		return api.VERDICT_COMPILATION_ERROR
	}

	if expected != actual {
		log.Printf("[%d] incorrect result on test::\nExpected: %s\nActual: %s\n", id, expected, actual)
		return api.VERDICT_FAILED
	}

	log.Printf("[%d] correct result on test:\nExpected: %s\nActual: %s\n", id, expected, actual)
	return api.VERDICT_PASSED
}

func TestSolution(tempFile string, problemId int) *api.SubmissionResult {
	test, err := inmemory.GetInstance().GetTestByID(problemId)
	if err != nil {
		return &api.SubmissionResult{
			Verdict: api.VERDICT_UNKNOWN_ERROR,
		}
	}
	if err := PrepareExe(test.SolutionFile, tempFile); err != nil {
		return &api.SubmissionResult{
			Verdict: api.VERDICT_COMPILATION_ERROR,
		}
	}
	defer CleanUp()
	for id, testCase := range test.TestCases {
		verdict := RunTestCase(test.SolutionFile, tempFile, id+1, testCase)
		if verdict != api.VERDICT_PASSED {
			return &api.SubmissionResult{
				Verdict:    verdict,
				FailedTest: int32(testCase.Id + 1),
			}
		}
	}
	return &api.SubmissionResult{
		Verdict: api.VERDICT_PASSED,
	}
}
