package handler

import "site/internal/services"

type UserHandlerOption func(u *UserHandler)
type SubmissionHandlerOption func(s *SubmissionHandler)

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