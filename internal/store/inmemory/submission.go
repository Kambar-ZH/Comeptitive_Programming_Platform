package inmemory

import (
	"context"
	"fmt"
	"site/internal/grpc/api"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubmissionRepo struct {
	data map[int32]*api.Submission
	api.UnimplementedSubmissionRepositoryServer
	mu   *sync.RWMutex
}

func (db *SubmissionRepo) All(ctx context.Context, empty *api.Pagination) (*api.SubmissionList, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := []*api.Submission{}
	for _, submission := range db.data {
		res = append(res, submission)
	}
	ans := api.SubmissionList{Submissions: res}

	return &ans, nil
}

func (db *SubmissionRepo) ById(ctx context.Context, submission *api.SubmissionRequestId) (*api.Submission, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if submission, ok := db.data[submission.Id]; ok {
		return submission, nil
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("submission with id %d does not exist", submission.Id))
}

func (db *SubmissionRepo) Create(ctx context.Context, submission *api.Submission) (*api.Submission, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[submission.Id] = submission
	return submission, nil
}

func (db *SubmissionRepo) Update(ctx context.Context, submission *api.Submission) (*api.Submission, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[submission.Id] = submission

	return submission, nil
}

func (db *SubmissionRepo) Delete(ctx context.Context, submission *api.SubmissionRequestId) (*api.Empty, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.data[submission.Id]; !ok {
		return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
	}
	delete(db.data, submission.Id)
	return &api.Empty{}, nil
}