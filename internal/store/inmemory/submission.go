package inmemory

import (
	"context"
	"fmt"
	"site/internal/datastruct"
	"site/internal/dto"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubmissionRepo struct {
	data map[int32]*datastruct.Submission
	mu   *sync.RWMutex
}


func (db *SubmissionRepo) All(ctx context.Context, query *dto.SubmissionFindAllRequest) ([]*datastruct.Submission, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := []*datastruct.Submission{}
	for _, submission := range db.data {
		res = append(res, submission)
	}

	return res, nil
}

func (db *SubmissionRepo) ById(ctx context.Context, id int) (*datastruct.Submission, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if submission, ok := db.data[int32(id)]; ok {
		return submission, nil
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("submission with id %d does not exist", id))
}

func (db *SubmissionRepo) Create(ctx context.Context, submission *datastruct.Submission) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[submission.Id] = submission
	return nil
}

func (db *SubmissionRepo) Update(ctx context.Context, submission *datastruct.Submission) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[submission.Id] = submission
	return nil
}

func (db *SubmissionRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	delete(db.data, int32(id))
	return nil
}