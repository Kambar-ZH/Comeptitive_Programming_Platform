package handler

import "site/internal/services"

type UserHandlerOption func(u *UserHandler)
type SubmissionHandlerOption func(s *SubmissionHandler)

func WithUserRepo(service services.UserService) UserHandlerOption {
	return func(u *UserHandler) {
		u.service = service
	}
}

func WithSubmissionRepo(service services.SubmissionService) SubmissionHandlerOption {
	return func(s *SubmissionHandler) {
		s.service = service
	}
}