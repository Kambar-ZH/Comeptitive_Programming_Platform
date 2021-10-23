package inmemory

import (
	"context"
	"fmt"
	"site/internal/models"
	"sync"
)

type SubmissionsRepo struct {
	data map[int]*models.Submission
	mu   *sync.RWMutex
}

func (db *SubmissionsRepo) Create(ctx context.Context, submission *models.Submission) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[submission.Id] = submission
	return nil
}

func (db *SubmissionsRepo) All(ctx context.Context) ([]*models.Submission, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	submissions := make([]*models.Submission, 0, len(db.data))
	for _, submission := range db.data {
		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (db *SubmissionsRepo) ById(ctx context.Context, id int) (*models.Submission, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	submission, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no submission with id %d", id)
	}
	return submission, nil
}

func (db *SubmissionsRepo) Update(ctx context.Context, submission *models.Submission) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[submission.Id] = submission
	return nil
}

func (db *SubmissionsRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
