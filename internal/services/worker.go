package services

import (
	"log"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/tools"
	"sync"
)

var (
	mu = &sync.Mutex{}
)

type Worker struct {
	command                string
	cleanFile              string
	mainSolution           string
	mainSolutionExe        string
	participantSolution    string
	participantSolutionExe string
}

func (worker *Worker) PrepareExe(solutionFile, tempFile string) error {
	if err := tools.CopyFile(worker.mainSolution, solutionFile); err != nil {
		return err
	}
	if err := tools.CopyFile(worker.participantSolution, tempFile); err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	if err := tools.BuildExe(worker.command); err != nil {
		log.Println("error on building exe")
		return err
	}
	return nil
}

func (worker *Worker) CleanUp() error {
	if err := tools.CopyFile(worker.participantSolution, worker.cleanFile); err != nil {
		return err
	}
	if err := tools.CopyFile(worker.mainSolution, worker.cleanFile); err != nil {
		return err
	}
	return nil
}
 
func (worker *Worker) RunTestCase(testCase datastruct.TestCase) (dto.Verdict, error) {
	expected, err := tools.MustExecuteFile(worker.mainSolutionExe, testCase)
	if err != nil {
		log.Println("error on executing main solution")
		return dto.UNKNOWN_ERROR, err
	}
	actual, err := tools.MustExecuteFile(worker.participantSolutionExe, testCase)
	if err != nil {
		log.Println("error on executing participant solution")
		return dto.COMPILATION_ERROR, err
	}

	if expected != actual {
		log.Printf("[%d] incorrect result on test:\nExpected: %s\nActual: %s\n", testCase.Id, expected, actual)
		return dto.FAILED, err
	}

	log.Printf("[%d] correct result on test:\nExpected: %s\nActual: %s\n", testCase.Id, expected, actual)
	return dto.PASSED, err
}