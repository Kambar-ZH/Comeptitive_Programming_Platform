package handler

import (
	"site/internal/services"

	"github.com/gorilla/sessions"
)

type UserHandlerOption func(u *UserHandler)
type SubmissionHandlerOption func(s *SubmissionHandler)
type UploadFileHandlerOption func(uf *UploadFileHandler)
type AuthHandlerOption func(s *AuthHandler)
type ContestHanderOption func(c *ContestHander)

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

func WithAuthService(service services.AuthService) AuthHandlerOption {
	return func(a *AuthHandler) {
		a.service = service
	}
}

func WithSessionStore(sessioStore sessions.Store) AuthHandlerOption {
	return func(a *AuthHandler) {
		a.sessionStore = sessioStore
	}
}

func WithContestService(service services.ContestService) ContestHanderOption {
	return func(c *ContestHander) {
		c.service = service
	}
}