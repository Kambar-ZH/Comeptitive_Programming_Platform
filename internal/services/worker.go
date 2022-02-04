package services

import (
	"fmt"
	"os"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/logger"
	"site/internal/tools"
	"site/test/inmemory"

	"github.com/hashicorp/go-uuid"
)

type Worker struct {
	mainDirectory    string
	mainSolutionFile string
	mainSolutionExec string
	userDirectory    string
	userSolutionFile string
	userSolutionExec string
}

func NewWorker() (*Worker, error) {
	uniqueId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	worker := &Worker{}

	worker.userDirectory = fmt.Sprintf("%s/user-%s", inmemory.MakeMe(), uniqueId)
	worker.userSolutionFile = fmt.Sprintf("%s/sol.go", worker.userDirectory)
	worker.userSolutionExec = fmt.Sprintf("%s/sol.exe", worker.userDirectory)

	worker.mainDirectory = fmt.Sprintf("%s/main-%s", inmemory.MakeMe(), uniqueId)
	worker.mainSolutionFile = fmt.Sprintf("%s/sol.go", worker.mainDirectory)
	worker.mainSolutionExec = fmt.Sprintf("%s/sol.exe", worker.mainDirectory)

	if err = worker.CreateDirs(); err != nil {
		return nil, err
	}
	return worker, nil
}

func (worker *Worker) CreateDirs() error {
	if err := os.Mkdir(worker.userDirectory, 0700); err != nil {
		return err
	}
	if err := os.Mkdir(worker.mainDirectory, 0700); err != nil {
		return err
	}
	return nil
}

func (worker *Worker) PrepareExe(solFile, tempFile, solExec string) error {
	if err := tools.CopyFile(tempFile, solFile); err != nil {
		return err
	}
	if err := tools.BuildExe(solExec, solFile); err != nil {
		logger.Logger.Sugar().Debugf("error on building %s", solExec)
		return err
	}
	return nil
}

func (worker *Worker) PrepareMainExe(solFile string) error {
	return worker.PrepareExe(worker.mainSolutionFile, solFile, worker.mainSolutionExec)
}

func (worker *Worker) PrepareUserExe(tempFile string) error {
	return worker.PrepareExe(worker.userSolutionFile, tempFile, worker.userSolutionExec)
}

func (worker *Worker) RemoveAll() error {
	if err := os.RemoveAll(worker.mainDirectory); err != nil {
		logger.Logger.Error(err.Error())
	}
	if err := os.RemoveAll(worker.userDirectory); err != nil {
		logger.Logger.Error(err.Error())
	}
	return nil
}

func (worker *Worker) RunTestCase(testCase datastruct.TestCase) (dto.Verdict, error) {
	expected, err := tools.MustExecuteFile(worker.mainSolutionExec, testCase)
	if err != nil {
		logger.Logger.Error("error on executing main solution")
		return dto.UNKNOWN_ERROR, err
	}
	actual, err := tools.MustExecuteFile(worker.userSolutionExec, testCase)
	if err != nil {
		logger.Logger.Error("error on executing participant solution")
		return dto.COMPILATION_ERROR, err
	}

	if !worker.Check(expected, actual) {
		logger.Logger.Sugar().Debugf("incorrect result on test [%d]:\nExpected: %s\nActual: %s", testCase.Id, expected, actual)
		return dto.FAILED, err
	}

	logger.Logger.Sugar().Debugf("correct result on test [%d]:\nExpected: %s\nActual: %s", testCase.Id, expected, actual)
	return dto.PASSED, err
}

func (worker *Worker) Check(expected, actual string) bool {
	// TODO: make smart validation using checker
	// If checker not found run for equlity

	return expected == actual
}
