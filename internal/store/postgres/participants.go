package postgres

import (
	"context"
	"site/internal/datastruct"
	"site/internal/dto"
	"site/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Participants() store.ParticipantRepository {
	if db.participants == nil {
		db.participants = NewParticipantRepository(db.conn)
	}
	return db.participants
}

type ParticipantRepository struct {
	conn *sqlx.DB
}

func NewParticipantRepository(conn *sqlx.DB) store.ParticipantRepository {
	return &ParticipantRepository{conn: conn}
}

func (p ParticipantRepository) FindAll(ctx context.Context, req *dto.ParticipantFindAllRequest) ([]*datastruct.Participant, error) {
	participants := make([]*datastruct.Participant, 0)
	if req.Filter != "" {
		if err := p.conn.Select(&participants,
			`SELECT participants.* 
				FROM participants, users 
				WHERE participants.contest_id = $1
				AND users.handle ILIKE $2
				OFFSET $3
				LIMIT $4`,
			req.ContestId, "%"+req.Filter+"%", req.Offset, req.Limit); err != nil {
			return nil, err
		}
		return participants, nil
	}
	if err := p.conn.Select(&participants,
		`SELECT * 
			FROM participants 
			WHERE participants.contest_id = $1
			OFFSET $2
			LIMIT $3`,
		req.ContestId, req.Offset, req.Limit); err != nil {
		return nil, err
	}
	for i := range participants {
		problemResults := make([]datastruct.ProblemResults, 0)
		if err := p.conn.Select(&problemResults, 
		`SELECT *
			FROM problem_results
			WHERE contest_id = $1
			AND user_id = $2`,
		participants[i].ContestId, participants[i].UserId); err != nil {
			return nil, err
		}
		participants[i].ProblemResults = problemResults
	}
	return participants, nil
}

func (p ParticipantRepository) FindFriends(ctx context.Context, req *dto.ParticipantFindFriendsRequest) ([]*datastruct.Participant, error) {
	participants := make([]*datastruct.Participant, 0)
	if req.Filter != "" {
		if err := p.conn.Select(&participants,
			`SELECT participants.* 
				FROM participants, users, user_friends
				WHERE participants.contest_id = $1
				AND user_friends.user_id = participants.user_id
				AND users.id = user_friends.friend_id
				AND users.handle ILIKE $2
				OFFSET $3
				LIMIT $4`,
			req.ContestId, "%"+req.Filter+"%", req.Offset, req.Limit); err != nil {
			return nil, err
		}
		return participants, nil
	}
	if err := p.conn.Select(&participants,
		`SELECT participants.* 
			FROM participants, users, user_friends 
			WHERE participants.contest_id = $1
			AND user_friends.user_id = participants.user_id 
			AND users.id = user_friends.friend_id
			OFFSET $2
			LIMIT $3`,
		req.ContestId, req.Offset, req.Limit); err != nil {
		return nil, err
	}
	for i := range participants {
		problemResults := make([]datastruct.ProblemResults, 0)
		if err := p.conn.Select(&problemResults, 
		`SELECT *
			FROM problem_results
			WHERE contest_id = $1
			AND user_id = $2`,
		req.ContestId,
		participants[i].UserId); err != nil {
			return nil, err
		}
		participants[i].ProblemResults = problemResults
	}
	return participants, nil
}

func (p ParticipantRepository) GetByUserId(ctx context.Context, req *dto.ParticipantGetByUserIdRequest) (*datastruct.Participant, error) {
	participant := new(datastruct.Participant)
	if err := p.conn.Get(participant,
		`SELECT * 
			FROM participants 
			WHERE user_id = $1
			AND contest_id = $2`,
		req.UserId, req.ContestId); err != nil {
		return nil, err
	}
	return participant, nil
} 

func (p ParticipantRepository) Create(ctx context.Context, participant *datastruct.Participant) error {
	_, err := p.conn.Exec(
		`INSERT INTO 
			participants (user_id, contest_id, participant_type, room) 
			VALUES ($1, $2, $3, $4)`,
			participant.UserId, participant.ContestId, participant.ParticipantType, participant.Room,
		)
	if err != nil {
		return err
	}
	return nil
}