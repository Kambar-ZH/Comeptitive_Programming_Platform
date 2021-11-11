package handler

import "site/internal/services"

type UserHandlerOption func(u *UserHandler)
type SubmissionHandlerOption func(s *SubmissionHandler)
type UploadFileHandlerOption func(s *UploadFileHandler)


func WithUserService(service services.UserService) UserHandlerOption {
	return func(u *UserHandler) {
		u.service = service
	}
}

func WithSubmissionService(service services.SubmissionService) SubmissionHandlerOption {
	return func(s *SubmissionHandler) {
		s.service = service
	}
}

func WithUploadFileService(service services.UploadFileService) UploadFileHandlerOption {
	return func(uf *UploadFileHandler) {
		uf.service = service
	}
}