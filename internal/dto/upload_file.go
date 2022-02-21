package dto

import (
	"mime/multipart"
	"site/internal/consts"
	"site/internal/datastruct"
)

type (
	UploadFileRequest struct {
		File      multipart.File
		FileName  string
		ProblemId int
		ContestId int
	}
	UploadFileResponse struct {
		Submission *datastruct.Submission
		Error      error
	}
	RunTestCasesRequest struct {
		ParticipantSolutionFilePath string
		ProblemId                   int
	}
	RunTestCasesResponse struct {
		Verdict    consts.Verdict
		FailedTest int32
	}
)
