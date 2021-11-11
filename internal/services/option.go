package services

import "site/internal/store"

type UserServiceOption func(u *UserServiceImpl)
type SubmissionServiceOption func(s *SubmissionServiceImpl)
type AuthServiceOption func(s *AuthServiceImpl)
type UploadFileServiceOption func(s *UploadFileServiceImpl)


func WithUserRepo(repo store.UserRepository) UserServiceOption {
	return func(u *UserServiceImpl) {
		u.repo = repo
	}
}

func WithSubmissionRepo(repo store.SubmissionRepository) SubmissionServiceOption {
	return func(s *SubmissionServiceImpl) {
		s.repo = repo
	}
}