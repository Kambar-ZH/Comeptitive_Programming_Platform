package services

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"site/internal/models"
	"site/test/inmemory"
	"time"
)

type UploadFileService interface {
	TestSolution(file string, problemId int) *models.SubmissionResult
	SaveInmemory(file multipart.File) (string, error)
}

type UploadFileServiceImpl struct {
	
}

func NewUploadFileService(opts ...UploadFileServiceOption) UploadFileService {
	svc := &UploadFileServiceImpl{}
	for _, v := range(opts) {
		v(svc)
	}
	return svc
}

func (c UploadFileServiceImpl) SaveInmemory(file multipart.File) (string, error) {
	tempFile, err := ioutil.TempFile(inmemory.GetInstance().GetTempSolutions(), "upload-*")
	if err != nil {
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileBytes)
	defer tempFile.Close()
	return tempFile.Name(), nil
}

func (c UploadFileServiceImpl) TestSolution(file string, problemId int) *models.SubmissionResult {
	test, err := inmemory.GetInstance().GetTestByID(problemId)
	if err != nil {
		return &models.SubmissionResult{
			Verdict: models.UNKNOWN_ERROR,
		}
	}
	if err := PrepareExe(test.SolutionFile, file); err != nil {
		return &models.SubmissionResult{
			Verdict: models.COMPILATION_ERROR,
		}
	}
	defer CleanUp()
	for id, testCase := range test.TestCases {
		verdict := RunTestCase(test.SolutionFile, file, id+1, testCase)
		if verdict != models.PASSED {
			return &models.SubmissionResult{
				Verdict:    verdict,
				FailedTest: testCase.Id + 1,
			}
		}
	}
	return &models.SubmissionResult {
		Verdict: models.PASSED,
	}
}

func CopyFile(inFile, fromFile string) error {
	fo, err := os.Create(inFile)
	if err != nil {
		return err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			return
		}
	}()

	data, err := os.ReadFile(fromFile)
	if err != nil {
		return err
	}
	fo.Write(data)
	return nil
}

func BuildExe() error {
	cmd := exec.Command("make")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ExecuteFile(file string, tCase models.TestCase) (string, error) {
	cmd := exec.Command(file)

	stdin, err := cmd.StdinPipe()
	defer func() {
		if err := stdin.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	defer func() {
		if err := stdout.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	input, err := os.ReadFile(tCase.InputFile)
	if err != nil {
		return "", err
	}

	_, err = stdin.Write(input)
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func MustExecuteFile(file string, tCase models.TestCase) (string, error) {
	ch := make(chan string)
	errorCh := make(chan error)

	go func() {
		verdict, err := ExecuteFile(file, tCase); 
		if err == nil {
			ch <- verdict
			return
		}
		errorCh <- err
	}()

	select {
	case verdict := <-ch:
		return verdict, nil
	case err := <-errorCh:
		return "", err
	case <-time.After(3 * time.Second):
		return "", fmt.Errorf("solution timeout")
	}
}

func PrepareExe(solutionFile, tempFile string) error {
	if err := CopyFile(inmemory.GetInstance().GetMainSolution(), solutionFile); err != nil {
		return err
	}
	if err := CopyFile(inmemory.GetInstance().GetParticipantSolution(), tempFile); err != nil {
		return err
	}
	if err := BuildExe(); err != nil {
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

func RunTestCase(solutionFile, tempFile string, id int, testCase models.TestCase) models.Verdict {
	expected, err := MustExecuteFile(inmemory.GetInstance().GetMainSolutionExe(), testCase)
	if err != nil {
		return models.UNKNOWN_ERROR
	}
	actual, err := MustExecuteFile(inmemory.GetInstance().GetParticipantSolutionExe(), testCase)
	if err != nil {
		return models.COMPILATION_ERROR
	}

	if expected != actual {
		log.Printf("[%d] incorrect result on test::\nExpected: %s\nActual: %s\n", id, expected, actual)
		return models.FAILED
	}

	log.Printf("[%d] correct result on test:\nExpected: %s\nActual: %s\n", id, expected, actual)
	return models.PASSED
}