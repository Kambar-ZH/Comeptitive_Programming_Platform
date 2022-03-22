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
		ProblemId int32
		ContestId int32
	}
	UploadFileResponse struct {
		Submission *datastruct.Submission
		Error      error
	}
	RunTestCasesRequest struct {
		ParticipantSolutionFilePath string
		ProblemId                   int32
	}
	RunTestCasesResponse struct {
		Verdict    consts.Verdict
		FailedTest int32
	}
)
