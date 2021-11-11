package services

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"site/internal/datastruct"
	"site/internal/middleware"
	"site/internal/models"
	"site/internal/store"
	"site/test/inmemory"
	"time"
)

type UploadFileService interface {
	TestSolution(file string, problemId int) (*models.SubmissionResult, error)
	SaveInmemory(file multipart.File) (string, error)
	Create(ctx context.Context, submission *datastruct.Submission) error
}

type UploadFileServiceImpl struct {
	repo store.SubmissionRepository
}

func (u UploadFileServiceImpl) Create(ctx context.Context, submission *datastruct.Submission) error {
	// TODO: ASSIGN USER TO SUBMISSION
	_, err := u.repo.ById(ctx, int(submission.Id))
	if err != nil {
		return err
	}
	user := middleware.UserFromCtx(ctx)
	submission.AuthorHandle = user.Handle
	return u.repo.Create(ctx, submission)
}

func NewUploadFileService(opts ...UploadFileServiceOption) UploadFileService {
	svc := &UploadFileServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (u UploadFileServiceImpl) SaveInmemory(file multipart.File) (string, error) {
	tempFile, err := ioutil.TempFile(inmemory.GetInstance().GetTempSolutions(), "upload-*")
	if err != nil {
		log.Println("couldn't save in memory")
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("couldn't read file")
		return "", err
	}
	tempFile.Write(fileBytes)
	defer tempFile.Close()
	return tempFile.Name(), nil
}

func (u UploadFileServiceImpl) TestSolution(file string, problemId int) (*models.SubmissionResult, error) {
	test, err := inmemory.GetInstance().GetTestByID(problemId)
	if err != nil {
		return &models.SubmissionResult{
			Verdict: models.UNKNOWN_ERROR,
		}, err
	}
	if err := PrepareExe(test.SolutionFile, file); err != nil {
		return &models.SubmissionResult{
			Verdict: models.COMPILATION_ERROR,
		}, err
	}
	defer CleanUp()
	for id, testCase := range test.TestCases {
		verdict, err := RunTestCase(test.SolutionFile, file, id+1, testCase)
		if err != nil {
			return nil, err
		}
		if verdict != models.PASSED {
			return &models.SubmissionResult{
				Verdict:    verdict,
				FailedTest: testCase.Id + 1,
			}, nil
		}
	}
	return &models.SubmissionResult{
		Verdict: models.PASSED,
	}, nil
}

func CopyFile(inFile, fromFile string) error {
	fo, err := os.Create(inFile)
	if err != nil {
		log.Println("couldn't create file")
		return err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Println("couldn't close infile")
			return
		}
	}()

	data, err := os.ReadFile(fromFile)
	if err != nil {
		log.Println("couldn't read file")
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
		log.Println("couldn't build exe")
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
		log.Println("error on stdinpipe")
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	defer func() {
		if err := stdout.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		log.Println("error on stdoutpipe")
		return "", err
	}

	if err := cmd.Start(); err != nil {
		log.Println("error on start cmd")
		return "", err
	}

	input, err := os.ReadFile(tCase.InputFile)
	if err != nil {
		log.Println("error on readfile to cmd")
		return "", err
	}

	_, err = stdin.Write(input)
	if err != nil {
		log.Println("error on write input")
		return "", err
	}

	result, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("error read from stdout")
		return "", err
	}
	return string(result), nil
}

func MustExecuteFile(file string, tCase models.TestCase) (string, error) {
	ch := make(chan string)
	errorCh := make(chan error)

	go func() {
		verdict, err := ExecuteFile(file, tCase)
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
		log.Println("error on copying to main solution file")
		return err
	}
	if err := CopyFile(inmemory.GetInstance().GetParticipantSolution(), tempFile); err != nil {
		log.Println("error on copying to participant solution file")
		return err
	}
	if err := BuildExe(); err != nil {
		log.Println("error on building exe")
		return err
	}
	return nil
}

func CleanUp() error {
	if err := CopyFile(inmemory.GetInstance().GetParticipantSolution(), inmemory.GetInstance().GetCleanFile()); err != nil {
		log.Println("error on cleaning participant solution")
		return err
	}
	if err := CopyFile(inmemory.GetInstance().GetMainSolution(), inmemory.GetInstance().GetCleanFile()); err != nil {
		log.Println("error on cleaning main solution")
		return err
	}
	return nil
}

func RunTestCase(solutionFile, tempFile string, id int, testCase models.TestCase) (models.Verdict, error) {
	expected, err := MustExecuteFile(inmemory.GetInstance().GetMainSolutionExe(), testCase)
	if err != nil {
		log.Println("error on executing main solution")
		return models.UNKNOWN_ERROR, err
	}
	actual, err := MustExecuteFile(inmemory.GetInstance().GetParticipantSolutionExe(), testCase)
	if err != nil {
		log.Println("error on executing participant solution")
		return models.COMPILATION_ERROR, err
	}

	if expected != actual {
		log.Printf("[%d] incorrect result on test::\nExpected: %s\nActual: %s\n", id, expected, actual)
		return models.FAILED, err
	}

	log.Printf("[%d] correct result on test:\nExpected: %s\nActual: %s\n", id, expected, actual)
	return models.PASSED, err
}