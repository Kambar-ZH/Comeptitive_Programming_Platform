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
}

func NewUploadFileService(opts ...UploadFileServiceOption) UploadFileService {
	svc := &UploadFileServiceImpl{}
	for _, v := range opts {
		v(svc)
	}
	return svc
}

func (u UploadFileServiceImpl) Create(ctx context.Context, submission *datastruct.Submission) error {
	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}
	submission.AuthorHandle = user.Handle
	return u.store.Submissions().Create(ctx, submission)
}

func (u UploadFileServiceImpl) SaveInmemory(dir string, file multipart.File) (string, error) {
	tempFile, err := ioutil.TempFile(dir, "upload-*")
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

func (u UploadFileServiceImpl) UploadFile(ctx context.Context, req *dto.UploadFileRequest) (*datastruct.Submission, error) {
	filePath, err := u.SaveInmemory(inmemory.TempSolutions(), req.File)
	if err != nil {
		return nil, err
	}
	res, err := u.RunTestCases(ctx, &dto.RunTestCasesRequest{
		FilePath:  filePath,
		ProblemId: req.ProblemId,
	})
	if err != nil {
		return nil, err
	}
	submission := &datastruct.Submission{
		Verdict:    string(res.Verdict),
		FailedTest: res.FailedTest,
		ContestId:  int32(req.ContestId),
		ProblemId: int32(req.ProblemId),
	}
	if err = u.Create(ctx, submission); err != nil {
		log.Println(err)
	}
	return submission, nil
}

func (u UploadFileServiceImpl) RunTestCases(ctx context.Context, req *dto.RunTestCasesRequest) (*dto.RunTestCasesResponse, error) {
	validator, err := u.store.Validators().ByProblemId(ctx, req.ProblemId)
	if err != nil {
		return &dto.RunTestCasesResponse{Verdict: dto.UNKNOWN_ERROR}, err
	}

	worker, err := NewWorker()
	if err != nil {
		return &dto.RunTestCasesResponse{Verdict: dto.UNKNOWN_ERROR}, err
	}
	defer worker.CleanUp()

	if err := worker.PrepareExe(validator.SolutionFilePath, req.FilePath); err != nil {
		return &dto.RunTestCasesResponse{Verdict: dto.COMPILATION_ERROR}, err
	}

	for _, testCase := range validator.TestCases {
		verdict, err := worker.RunTestCase(testCase)
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
