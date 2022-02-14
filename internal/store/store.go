package store

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
)

type Store interface {
	Connect(url string) error
	Close() error

	Users() UserRepository
	Submissions() SubmissionRepository
	Validators() ValidatorRepository
	TestCases() TestCaseRepository
	Contests() ContestRepository
	Problems() ProblemRepository
	Participants() ParticipantRepository
	ProblemResults() ProblemResultsRepository
	Tags() TagRepository
	UserFriends() UserFriendRepository
}

type UserRepository interface {
	FindAll(ctx context.Context, req *dto.UserFindAllRequest) ([]*datastruct.User, error)
	FindFriends(ctx context.Context, req *dto.UserFindFriendsRequest) ([]*datastruct.User, error)
	GetByEmail(ctx context.Context, email string) (*datastruct.User, error)
	GetByHandle(ctx context.Context, handle string) (*datastruct.User, error)
	Create(ctx context.Context, user *datastruct.User) error
	Update(ctx context.Context, user *datastruct.User) error
	Delete(ctx context.Context, handle string) error
}

type SubmissionRepository interface {
	FindAll(ctx context.Context, req *dto.SubmissionFindAllRequest) ([]*datastruct.Submission, error)
	GetById(ctx context.Context, id int) (*datastruct.Submission, error)
	Create(ctx context.Context, submission *datastruct.Submission) error
	Update(ctx context.Context, submission *datastruct.Submission) error
	Delete(ctx context.Context, id int) error
}

type ValidatorRepository interface {
	GetByProblemId(ctx context.Context, problemId int) (*datastruct.Validator, error)
}

type TestCaseRepository interface {
	GetByProblemId(ctx context.Context, problemId int) ([]*datastruct.TestCase, error)
}

type ContestRepository interface {
	FindAll(ctx context.Context, req *dto.ContestFindAllRequest) ([]*datastruct.Contest, error)
	FindByTimeInterval(ctx context.Context, req *dto.ContestFindByTimeInterval) ([]*datastruct.Contest, error)
	GetById(ctx context.Context, id int) (*datastruct.Contest, error)
	Create(ctx context.Context, contest *datastruct.Contest) error
	Update(ctx context.Context, contest *datastruct.Contest) error
	Delete(ctx context.Context, id int) error
}

type ProblemRepository interface {
	Problemset(ctx context.Context, req *dto.ProblemsetRequest) ([]*datastruct.Problem, error)
	FindAll(ctx context.Context, req *dto.ProblemFindAllRequest) ([]*datastruct.Problem, error)
	GetById(ctx context.Context, id int) (*datastruct.Problem, error)
	Create(ctx context.Context, problem *datastruct.Problem) error
	Update(ctx context.Context, problem *datastruct.Problem) error
	Delete(ctx context.Context, id int) error
}

type ParticipantRepository interface {
	FindAll(ctx context.Context, req *dto.ParticipantFindAllRequest) ([]*datastruct.Participant, error)
	FindFriends(ctx context.Context, req *dto.ParticipantFindFriendsRequest) ([]*datastruct.Participant, error)
	GetByUserId(ctx context.Context, req *dto.ParticipantGetByUserIdRequest) (*datastruct.Participant, error)
	Create(ctx context.Context, participant *datastruct.Participant) error
}

type ProblemResultsRepository interface {
	Update(ctx context.Context, problemResults *datastruct.ProblemResult) error
	GetByProblemId(ctx context.Context, req *dto.ProblemResultGetByProblemIdRequest) (*datastruct.ProblemResult, error)
	Create(ctx context.Context, problemResults *datastruct.ProblemResult) error
}

type TagRepository interface {
	GetByProblemId(ctx context.Context, problemId int) ([]*datastruct.Tag, error)
}

type UserFriendRepository interface {
	Create(ctx context.Context, userFriend *datastruct.UserFriend) error
	Delete(ctx context.Context, userFriend *datastruct.UserFriend) error
}
