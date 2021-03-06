package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"
	"site/test/inmemory"

	"github.com/jmoiron/sqlx"
)

func (db *DB) TestCases() store.TestCaseRepository {
	if db.testCases == nil {
		db.testCases = NewTestCaseRepository(db.conn)
	}
	return db.testCases
}

type TestCaseRepository struct {
	conn *sqlx.DB
}

func NewTestCaseRepository(conn *sqlx.DB) store.TestCaseRepository {
	return &TestCaseRepository{conn: conn}
}

func (v TestCaseRepository) GetByProblemId(ctx context.Context, problemId int) ([]*datastruct.TestCase, error) {
	testCases := make([]*datastruct.TestCase, 0)
	err := v.conn.Select(&testCases,
		`SELECT * 
			FROM test_cases 
			WHERE problem_id = $1`,
		problemId)
	if err != nil {
		return nil, err
	}
	for i := range testCases {
		testCases[i].TestFile = inmemory.AbsPath(testCases[i].TestFile)
	}
	return testCases, nil
}
