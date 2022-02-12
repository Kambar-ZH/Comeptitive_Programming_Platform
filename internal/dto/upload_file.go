package dto

import (
	"mime/multipart"
	"site/internal/datastruct"
)

const (
	PASSED            Verdict = "Passed"
	FAILED            Verdict = "Failed Test"
	COMPILATION_ERROR Verdict = "Compilation Error"
	TIME_LIMIT_ERROR  Verdict = "Time Limit Error"
	UNKNOWN_ERROR     Verdict = "Unknown Error"
)

type (
	Verdict           string
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
		Verdict    Verdict
		FailedTest int32
	}
)
