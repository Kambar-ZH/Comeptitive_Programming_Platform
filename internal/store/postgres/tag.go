package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"site/internal/datastruct"
	"site/internal/store"
)

func (db *DB) Tags() store.TagRepository {
	if db.tags == nil {
		db.tags = NewTagRepository(db.conn)
	}
	return db.tags
}

type TagRepository struct {
	conn *sqlx.DB
}

func NewTagRepository(conn *sqlx.DB) store.TagRepository {
	return &TagRepository{conn: conn}
}

func (t TagRepository) GetByProblemId(ctx context.Context, problemId int) ([]*datastruct.Tag, error) {
	tags := make([]*datastruct.Tag, 0)
	if err := t.conn.Select(&tags,
		`SELECT tags.* 
			FROM tags, problems_tags 
			WHERE problems_tags.problem_id = $1 
			AND tags.id = problems_tags.tag_id`,
		problemId); err != nil {
		return nil, err
	}
	return tags, nil
}
