package store

import (
	"site/grpc/api"
)

type Store interface {
	Users() UserRepository
	Submissions() SubmissionRepository
}

type UserRepository interface {
	api.UserRepositoryServer
}

type SubmissionRepository interface {
	api.SubmissionRepositoryServer
}