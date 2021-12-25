package services

import (
	"fmt"
	"log"
	"os"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/tools"
	"site/test/inmemory"

	"github.com/hashicorp/go-uuid"
)

type Worker struct {
	userDirectory    string
	mainDirectory    string
	mainSolutionFile string
	userSolutionFile string
	mainSolutionExec string
	userSolutionExec string
}

func NewWorker() (*Worker, error) {
	worker := &Worker{}
	uniqueId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	worker.userDirectory = fmt.Sprintf("%s/user-%s", inmemory.MakeMe(), uniqueId)
	err = os.Mkdir(worker.userDirectory, 0700)
	if err != nil {
		return nil, err
	}
	worker.userSolutionFile = fmt.Sprintf("%s/sol.go", worker.userDirectory)
	worker.userSolutionExec = fmt.Sprintf("%s/sol.exe", worker.userDirectory)

	worker.mainDirectory = fmt.Sprintf("%s/main-%s", inmemory.MakeMe(), uniqueId)
	err = os.Mkdir(worker.mainDirectory, 0700)
	if err != nil {
		return nil, err
	}
	worker.mainSolutionFile = fmt.Sprintf("%s/sol.go", worker.mainDirectory)
	worker.mainSolutionExec = fmt.Sprintf("%s/sol.exe", worker.mainDirectory)
	return worker, nil
}

func (worker *Worker) PrepareExe(solutionFile, tempFile string) error {
	if err := tools.CopyFile(worker.mainSolutionFile, solutionFile); err != nil {
		return err
	}
	if err := tools.CopyFile(worker.userSolutionFile, tempFile); err != nil {
		return err
	}
	if err := tools.BuildExe(worker.mainSolutionExec, worker.mainSolutionFile); err != nil {
		log.Println("error on building main exe")
		return err
	}
	if err := tools.BuildExe(worker.userSolutionExec, worker.userSolutionFile); err != nil {
		log.Println("error on building user exe")
		return err
	}
	return nil
}

func (worker *Worker) CleanUp() error {
	err := os.RemoveAll(worker.mainDirectory)
	if err != nil {
		log.Println(err)
	}
	err = os.RemoveAll(worker.userDirectory)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (worker *Worker) RunTestCase(testCase datastruct.TestCase) (dto.Verdict, error) {
	expected, err := tools.MustExecuteFile(worker.mainSolutionExec, testCase)
	if err != nil {
		log.Println("error on executing main solution")
		return dto.UNKNOWN_ERROR, err
	}
	actual, err := tools.MustExecuteFile(worker.userSolutionExec, testCase)
	if err != nil {
		log.Println("error on executing participant solution")
		return dto.COMPILATION_ERROR, err
	}

	if expected != actual {
		log.Printf("incorrect result on test [%d]:\nExpected: %s\nActual: %s\n", testCase.Id, expected, actual)
		return dto.FAILED, err
	}

	log.Printf("correct result on test [%d]:\nExpected: %s\nActual: %s\n", testCase.Id, expected, actual)
	return dto.PASSED, err
}
