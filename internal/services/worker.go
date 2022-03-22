package services

import (
	"site/internal/consts"
	"site/internal/datastruct"
	"site/internal/logger"
	"site/internal/tools"
)

func RunTestCase(testCase *datastruct.TestCase, userSolutionFilePath, authorSolutionFilePath string) (consts.Verdict, error) {
	expected, err := tools.MustExecuteFile(authorSolutionFilePath, testCase)
	if err != nil {
		logger.Logger.Error("error on executing author solution")
		return consts.UNKNOWN_ERROR, err
	}
	actual, err := tools.MustExecuteFile(userSolutionFilePath, testCase)
	if err != nil {
		logger.Logger.Error("error on executing participant solution")
		return consts.COMPILATION_ERROR, err
	}

	if !check(expected, actual) {
		logger.Logger.Sugar().Debugf("incorrect result on test [%d]:\nExpected: %s\nActual: %s", testCase.Id, expected, actual)
		return consts.FAILED, err
	}

	logger.Logger.Sugar().Debugf("correct result on test [%d]:\nExpected: %s\nActual: %s", testCase.Id, expected, actual)
	return consts.PRETESTS_PASSED, err
}

func check(expected, actual string) bool {
	// TODO: make smart validation using checker
	// If checker not found run for equlity

	return expected == actual
}
