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
	"site/test/inmemory"
	"strings"
)

type UploadFileService interface {
	UploadFile(ctx context.Context, req *dto.UploadFileRequest) (*datastruct.Submission, error)
}

type UploadFileServiceImpl struct {
	store store.Store

	tasks   chan *Task
	workers []*Worker
}

func (u UploadFileServiceImpl) RunPool() {
	for _, worker := range u.workers {
		go u.RunWorker(worker)
	}
}

func NewUploadFileService(opts ...UploadFileServiceOption) UploadFileService {
	svc := &UploadFileServiceImpl{
		tasks: make(chan *Task, 100),
		workers: []*Worker{
			{
				command:                "all",
				cleanFile:              inmemory.AbsolutePath("cmd/myapp/main_solution/clean.txt"),
				mainSolution:           inmemory.AbsolutePath("cmd/myapp/main_solution/main_solution.go"),
				mainSolutionExe:        inmemory.AbsolutePath("cmd/myapp/main_solution/main_solution.exe"),
				participantSolution:    inmemory.AbsolutePath("cmd/myapp/participant_solution/participant_solution.go"),
				participantSolutionExe: inmemory.AbsolutePath("cmd/myapp/participant_solution/participant_solution.exe"),
			},
		},
	}
	for _, v := range opts {
		v(svc)
	}
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

func (u UploadFileServiceImpl) RunWorker(worker *Worker) {
	for task := range u.tasks {
		filePath, err := u.SaveInmemory(task.req.File)
		if err != nil {
			task.out <- &dto.UploadFileResponse{Submission: nil, Error: err}
			return
		}
		log.Println(filePath)

		res, err := u.RunTestCases(task.ctx, worker, &dto.RunTestCasesRequest{
			FilePath:  filePath,
			ProblemId: task.req.ProblemId,
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
	u.tasks <- &Task{
		req: req,
		ctx: ctx,
		out: out,
	}
	result := <-out
	return result.Submission, result.Error
}

func (u UploadFileServiceImpl) RunTestCases(ctx context.Context, worker *Worker, req *dto.RunTestCasesRequest) (*dto.RunTestCasesResponse, error) {
	validator, err := u.store.Validators().ByProblemId(ctx, req.ProblemId)
	if err != nil {
		return &dto.RunTestCasesResponse{Verdict: dto.UNKNOWN_ERROR}, err
	}

	if err := worker.PrepareExe(validator.SolutionFilePath, req.FilePath); err != nil {
		return &dto.RunTestCasesResponse{Verdict: dto.COMPILATION_ERROR}, err
	}

	defer worker.CleanUp()
	for id, testCase := range validator.TestCases {
		verdict, err := worker.RunTestCase(validator.SolutionFilePath, req.FilePath, id+1, testCase)
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
