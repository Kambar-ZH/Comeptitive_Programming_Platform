package services

import (
	"site/internal/store"
	"site/internal/cache"
)

type UserServiceOption func(u *UserServiceImpl)
type SubmissionServiceOption func(s *SubmissionServiceImpl)
type UploadFileServiceOption func(s *UploadFileServiceImpl)
type AuthServiceOption func(s *AuthServiceImpl)

func UserServiceWithStore(store store.Store) UserServiceOption {
	return func(u *UserServiceImpl) {
		u.store = store
	}
}

func SubmissionServiceWithCache(cache cache.Cache) SubmissionServiceOption {
	return func(s *SubmissionServiceImpl) {
		s.cache = cache
	}
}

func SubmissionServiceWithStore(store store.Store) SubmissionServiceOption {
	return func(s *SubmissionServiceImpl) {
		s.store = store
	}
}

func UploadFileServiceWithStore(store store.Store) UploadFileServiceOption {
	return func(uf *UploadFileServiceImpl) {
		uf.store = store
	}
}

func AuthServiceWithStore(store store.Store) AuthServiceOption {
	return func(a *AuthServiceImpl) {
		a.store = store
	}
}
