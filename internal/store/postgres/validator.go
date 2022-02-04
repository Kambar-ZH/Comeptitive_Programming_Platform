package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/store"

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

func (v ValidatorRepository) ByProblemId(ctx context.Context, problemId int) (*datastruct.Validator, error) {
	validator := new(datastruct.Validator)
	err := v.conn.Get(validator, 
		`SELECT * 
			FROM validators 
			WHERE problem_id = $1`, 
		problemId)
	if err != nil {
		return nil, err
	}
	testCases := make([]datastruct.TestCase, 0)
	err = v.conn.Select(&testCases, 
		`SELECT * 
			FROM test_cases 
			WHERE problem_id = $1`, 
		problemId)
	if err != nil {
		return nil, err
	}
	validator.TestCases = testCases
	return validator, nil
}