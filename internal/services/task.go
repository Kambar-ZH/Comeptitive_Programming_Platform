package services

import (
	"context"
	"site/internal/dto"
)

type Task struct {
	req *dto.UploadFileRequest

	ctx context.Context
	out chan *dto.UploadFileResponse
}