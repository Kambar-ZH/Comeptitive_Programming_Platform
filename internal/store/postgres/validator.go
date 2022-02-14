package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"
	"site/test/inmemory"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Validators() store.ValidatorRepository {
	if db.validators == nil {
		db.validators = NewValidatorRepository(db.conn)
	}
	return db.validators
}

type ValidatorRepository struct {
	conn *sqlx.DB
}

func NewValidatorRepository(conn *sqlx.DB) store.ValidatorRepository {
	return &ValidatorRepository{conn: conn}
}

func (v ValidatorRepository) ValidatorWithTestCases(ctx context.Context, validator *datastruct.Validator) (*datastruct.Validator, error) {
	testCases := make([]*datastruct.TestCase, 0)
	err := v.conn.Select(&testCases,
		`SELECT * 
			FROM test_cases 
			WHERE problem_id = $1`,
		validator.ProblemId)
	if err != nil {
		return nil, err
	}
	for i := range testCases {
		testCases[i].TestFile = inmemory.AbsPath(testCases[i].TestFile)
	}
	validator.TestCases = testCases
	return validator, nil
}

func (v ValidatorRepository) GetByProblemId(ctx context.Context, problemId int) (*datastruct.Validator, error) {
	validator := new(datastruct.Validator)
	err := v.conn.Get(validator,
		`SELECT * 
			FROM validators 
			WHERE problem_id = $1`,
		problemId)
	if err != nil {
		return nil, err
	}

	// Get Abs Path to validator and testcases
	validator.AuthorSolutionFilePath = inmemory.AbsPath(validator.AuthorSolutionFilePath)
	validator, err = v.ValidatorWithTestCases(ctx, validator)
	if err != nil {
		return nil, err
	}
	return validator, nil
}
