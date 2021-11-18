package services

import (
	"context"
	"io/ioutil"
	"log"
	"mime/multipart"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/middleware"
	"site/internal/store"
	"site/internal/tools"
	"site/test/inmemory"
	"strings"
)

type UploadFileService interface {
	UploadFile(ctx context.Context, req *dto.UploadFileRequest) (*datastruct.Submission, error)
}

type Task struct {
	req *dto.UploadFileRequest
	ctx context.Context
	out chan *dto.UploadFileResponse
}

type Pool struct {
	tasks    chan *Task
	commands []string
}

type UploadFileServiceImpl struct {
	store store.Store
	pool  *Pool
}

func (u UploadFileServiceImpl) RunPool() {
	for _, command := range u.pool.commands {
		go u.Worker(command)
	}
}	

func NewUploadFileService(opts ...UploadFileServiceOption) UploadFileService {
	svc := &UploadFileServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	svc.pool = &Pool{
		tasks:    make(chan *Task, 100),
		commands: []string{"all"},
	}
	svc.RunPool()
	return svc
}

func (u UploadFileServiceImpl) Create(ctx context.Context, submission *datastruct.Submission) error {
	// TODO: ASSIGN USER TO SUBMISSION
	_, err := u.store.Submissions().ById(ctx, int(submission.Id))
	if err != nil {
		return err
	}
	user := middleware.UserFromCtx(ctx)
	submission.AuthorHandle = user.Handle
	return u.store.Submissions().Create(ctx, submission)
}

func (u UploadFileServiceImpl) SaveInmemory(file multipart.File) (string, error) {
	tempFile, err := ioutil.TempFile(inmemory.TempSolutions(), "upload-*")
	if err != nil {
		return "", err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	tempFile.Write(fileBytes)
	defer tempFile.Close()
	filePath := strings.ReplaceAll(tempFile.Name(), "\\\\", "/")
	return filePath, nil
}

func (u *UploadFileServiceImpl) Worker(command string) {
	for task := range u.pool.tasks {
		filePath, err := u.SaveInmemory(task.req.File)
		if err != nil {
			task.out <- &dto.UploadFileResponse{Submission: nil, Error: err}
			return
		}
		log.Println(filePath)

		res, err := u.RunTestCases(task.ctx, &dto.RunTestCasesRequest{
			FilePath:  filePath,
			ProblemId: task.req.ProblemId,
			Command: command,
		})
		if err != nil {
			task.out <- &dto.UploadFileResponse{Submission: nil, Error: err}
			return
		}
		submission := &datastruct.Submission{
			Verdict:    string(res.Verdict),
			FailedTest: res.FailedTest,
		}
		u.Create(task.ctx, submission)
		task.out <- &dto.UploadFileResponse{Submission: submission, Error: nil}
	}
}

func (u UploadFileServiceImpl) UploadFile(ctx context.Context, req *dto.UploadFileRequest) (*datastruct.Submission, error) {
	out := make(chan *dto.UploadFileResponse)
	u.pool.tasks <- &Task{
		req: req,
		ctx: ctx,
		out: out,
	}
	result := <-out
	return result.Submission, result.Error
}

func (u UploadFileServiceImpl) RunTestCases(ctx context.Context, req *dto.RunTestCasesRequest) (*dto.RunTestCasesResponse, error) {
	validator, err := u.store.Validators().ByProblemId(ctx, req.ProblemId)
	if err != nil {
		return &dto.RunTestCasesResponse{Verdict: dto.UNKNOWN_ERROR}, err
	}
	if err := PrepareExe(validator.SolutionFilePath, req.FilePath, req.Command); err != nil {
		return &dto.RunTestCasesResponse{Verdict: dto.COMPILATION_ERROR}, err
	}
	defer CleanUp()
	for id, testCase := range validator.TestCases {
		verdict, err := RunTestCase(validator.SolutionFilePath, req.FilePath, id+1, testCase)
		if err != nil {
			return nil, err
		}
		if verdict != dto.PASSED {
			return &dto.RunTestCasesResponse{
				Verdict:    verdict,
				FailedTest: testCase.Id + 1,
			}, nil
		}
	}
	return &dto.RunTestCasesResponse{Verdict: dto.PASSED}, nil
}

func PrepareExe(solutionFile, tempFile, command string) error {
	if err := tools.CopyFile(inmemory.MainSolution(), solutionFile); err != nil {
		log.Println("error on copying to main solution file")
		return err
	}
	if err := tools.CopyFile(inmemory.ParticipantSolution(), tempFile); err != nil {
		log.Println("error on copying to participant solution file")
		return err
	}
	if err := tools.BuildExe(command); err != nil {
		log.Println("error on building exe")
		return err
	}
	return nil
}

func CleanUp() error {
	if err := tools.CopyFile(inmemory.ParticipantSolution(), inmemory.CleanFile()); err != nil {
		log.Println("error on cleaning participant solution")
		return err
	}
	if err := tools.CopyFile(inmemory.MainSolution(), inmemory.CleanFile()); err != nil {
		log.Println("error on cleaning main solution")
		return err
	}
	return nil
}

func RunTestCase(solutionFile, tempFile string, id int, testCase datastruct.TestCase) (dto.Verdict, error) {
	expected, err := tools.MustExecuteFile(inmemory.MainSolutionExe(), testCase)
	if err != nil {
		log.Println("error on executing main solution")
		return dto.UNKNOWN_ERROR, err
	}
	actual, err := tools.MustExecuteFile(inmemory.ParticipantSolutionExe(), testCase)
	if err != nil {
		log.Println("error on executing participant solution")
		return dto.COMPILATION_ERROR, err
	}

	if expected != actual {
		log.Printf("[%d] incorrect result on test::\nExpected: %s\nActual: %s\n", id, expected, actual)
		return dto.FAILED, err
	}

	log.Printf("[%d] correct result on test:\nExpected: %s\nActual: %s\n", id, expected, actual)
	return dto.PASSED, err
}