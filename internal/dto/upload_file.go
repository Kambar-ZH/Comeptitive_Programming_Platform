package dto

import "mime/multipart"

const (
	PASSED            Verdict = "Passed"
	FAILED            Verdict = "Failed Test"
	COMPILATION_ERROR Verdict = "Compilation Error"
	TIME_LIMIT_ERROR  Verdict = "Time Limit Error"
	UNKNOWN_ERROR     Verdict = "Unknown Error"
)

type (
	Verdict            string
	UploadFileRequest struct {
		File      multipart.File
		ProblemId int
	}
	RunTestCasesRequest struct {
		FilePath  string
		ProblemId int
	}
	RunTestCasesResponse struct {
		Verdict    Verdict
		FailedTest int32
	}
)
