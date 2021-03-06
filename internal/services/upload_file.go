package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"site/internal/consts"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/logger"
	"site/internal/middleware"
	"site/internal/store"
	"site/internal/tools"
	"site/test/inmemory"
	"time"
)

var (
	penaltyArrest = 50
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
	submission.UserId = user.Id
	if err := u.assignPoints(ctx, submission); err != nil {
		return err
	}
	return u.store.Submissions().Create(ctx, submission)
}

func (u UploadFileServiceImpl) assignPoints(ctx context.Context, submission *datastruct.Submission) error {
	contest, err := u.store.Contests().GetById(ctx, &dto.ContestGetByIdRequest{
		ContestId:    submission.ContestId,
		LanguageCode: consts.EN,
	})
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}
	if contest.Phase != consts.CODING.String() {
		logger.Logger.Debug("contest is not in coding phase")
		return nil
	}
	problem, err := u.store.Problems().GetById(ctx, &dto.ProblemGetByIdRequest{
		ProblemId:    submission.ProblemId,
		LanguageCode: consts.EN,
	})
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	points := problem.Points
	minutePassed := int32(time.Since(contest.StartDate) / time.Minute)
	contestDuration := int32(contest.EndDate.Sub(contest.StartDate) / time.Minute)

	if contestDuration != 0 {
		points -= minutePassed / contestDuration * points * 50
	}

	problemResults, err := u.store.ProblemResults().GetByProblemId(ctx, &dto.ProblemResultGetByProblemIdRequest{
		ContestId: submission.ContestId,
		ProblemId: submission.ProblemId,
		UserId:    submission.UserId,
	})
	if problemResults == nil {
		problemResults = &datastruct.ProblemResult{
			UserId:                       submission.UserId,
			ProblemId:                    submission.ProblemId,
			ContestId:                    submission.ContestId,
			Points:                       0,
			Penalty:                      0,
			LastSuccessfulSubmissionTime: time.Now(),
		}
		u.store.ProblemResults().Create(ctx, problemResults)
	}

	pointsBefore := problemResults.Points

	if submission.Verdict != string(consts.PRETESTS_PASSED) {
		problemResults.Penalty++
		points -= int32(penaltyArrest)
	}

	if err = u.store.ProblemResults().Update(ctx, problemResults); err != nil {
		logger.Logger.Error("couldn't update problem results")
		return err
	}

	user, ok := middleware.UserFromCtx(ctx)
	if !ok {
		return middleware.ErrNotAuthenticated
	}

	participant, err := u.store.Participants().GetById(ctx, &dto.ParticipantGetByIdRequest{
		ContestId: contest.Id,
		UserId:    user.Id,
	})

	participant.Points += points - pointsBefore

	if err = u.store.Participants().Update(ctx, participant); err != nil {
		logger.Logger.Error("couldn't update participant points")
		return err
	}

	return nil
}

func (u UploadFileServiceImpl) saveInMemory(dir string, file multipart.File, fileName string) (*os.File, error) {
	fileExtension := tools.ExtensionRegex.ReplaceAllString(fileName, "")
	tempFile, err := ioutil.TempFile(dir, fmt.Sprintf("*.%s", fileExtension))
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	tempFile.Write(fileBytes)
	return tempFile, nil
}

func (u UploadFileServiceImpl) UploadFile(ctx context.Context, req *dto.UploadFileRequest) (*datastruct.Submission, error) {
	// WARN: after saving file, closed
	file, err := u.saveInMemory(inmemory.TempSolutions(), req.File, req.FileName)
	if err != nil {
		return nil, err
	}
	filePath := file.Name()
	res, err := u.RunTestCases(ctx, &dto.RunTestCasesRequest{
		ParticipantSolutionFilePath: filePath,
		ProblemId:                   req.ProblemId,
	})
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}
	submission := &datastruct.Submission{
		Verdict:          res.Verdict.String(),
		FailedTest:       res.FailedTest,
		ContestId:        req.ContestId,
		ProblemId:        req.ProblemId,
		SolutionFilePath: filePath,
	}
	if err = u.Create(ctx, submission); err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}
	return submission, nil
}

func (u UploadFileServiceImpl) RunTestCases(ctx context.Context, req *dto.RunTestCasesRequest) (*dto.RunTestCasesResponse, error) {
	validator, err := u.store.Validators().GetByProblemId(ctx, int(req.ProblemId))
	if err != nil {
		logger.Logger.Error(err.Error())
		return &dto.RunTestCasesResponse{Verdict: consts.UNKNOWN_ERROR}, err
	}

	for _, testCase := range validator.TestCases {
		verdict, err := RunTestCase(testCase, req.ParticipantSolutionFilePath, validator.AuthorSolutionFilePath)
		if err != nil {
			return nil, err
		}
		if verdict != consts.PRETESTS_PASSED {
			return &dto.RunTestCasesResponse{
				Verdict:    verdict,
				FailedTest: testCase.Id,
			}, nil
		}
	}
	return &dto.RunTestCasesResponse{Verdict: consts.PRETESTS_PASSED}, nil
}
