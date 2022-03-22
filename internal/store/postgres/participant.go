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
		err := p.conn.Select(&participants, `
			SELECT 
				u.id     "user_id",
			   	p.room   "room", 
			    p.contest_id "contest_id",
			    p.participant_type "participant_type",
		       	u.handle "handle", 
		       	u.rating "rating",
			    p.points "points"
			FROM participants p
				JOIN users u on p.user_id = u.id
			WHERE p.contest_id = $1
				AND u.handle ILIKE $2
			OFFSET $3
			LIMIT $4`, req.ContestId, "%"+req.Filter+"%", req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}
		return participants, nil
	}
	err := p.conn.Select(&participants, `
		SELECT 
			u.id     "user_id",
			p.room   "room",
			p.contest_id "contest_id",
			p.participant_type "participant_type",
			u.handle "handle", 
			u.rating "rating",
			p.points "points"
		FROM participants p
			JOIN users u on p.user_id = u.id
		WHERE p.contest_id = $1
		OFFSET $2
		LIMIT $3`,
		req.ContestId, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (p ParticipantRepository) FindFriends(ctx context.Context, req *dto.ParticipantFindFriendsRequest) ([]*datastruct.Participant, error) {
	participants := make([]*datastruct.Participant, 0)
	if req.Filter != "" {
		err := p.conn.Select(&participants, `
			SELECT 
				u.id     "user_id",
			   	p.room   "room",
			    p.contest_id "contest_id",
			    p.participant_type "participant_type",
		       	u.handle "handle", 
		       	u.rating "rating",
			    p.points "points"
			FROM participants p
			    JOIN user_friends uf on uf.user_id = p.user_id
				JOIN users u on u.id = uf.friend_id
			WHERE p.contest_id = $1
				AND u.handle ILIKE $2
			OFFSET $3
			LIMIT $4`,
			req.ContestId, "%"+req.Filter+"%", req.Offset, req.Limit)
		if err != nil {
			return nil, err
		}
		return participants, nil
	}
	err := p.conn.Select(&participants, `
		SELECT 
			u.id     "user_id",
			p.room   "room",
			p.contest_id "contest_id",
			p.participant_type "participant_type",
			u.handle "handle", 
			u.rating "rating",
			p.points "points"
		FROM participants p
			JOIN user_friends uf on uf.user_id = p.user_id
			JOIN users u on u.id = uf.friend_id
		WHERE p.contest_id = $1
		OFFSET $2
		LIMIT $3`,
		req.ContestId, req.Offset, req.Limit)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (p ParticipantRepository) GetById(ctx context.Context, req *dto.ParticipantGetByIdRequest) (*datastruct.Participant, error) {
	participant := new(datastruct.Participant)
	err := p.conn.Get(participant, `
		SELECT  
			u.id     "user_id",
			p.room   "room",
			p.contest_id "contest_id",
			p.participant_type "participant_type",
			u.handle "handle", 
			u.rating "rating",
			p.points "points"
		FROM participants p
			JOIN users u on p.user_id = u.id
		WHERE p.contest_id = $1
			AND p.user_id = $2`,
		req.ContestId, req.UserId)
	if err != nil {
		return nil, err
	}
	return participant, nil
}

func (p ParticipantRepository) Create(ctx context.Context, participant *datastruct.Participant) error {
	_, err := p.conn.Exec(`
			INSERT INTO 
				participants (user_id, contest_id, participant_type, room) 
			VALUES ($1, $2, $3, $4)`,
		participant.UserId, participant.ContestId, participant.ParticipantType, participant.Room,
	)
	if err != nil {
		return err
	}
	return nil
}

func (u ParticipantRepository) Update(ctx context.Context, participant *datastruct.Participant) error {
	_, err := u.conn.Exec(`
		UPDATE participants 
			SET user_id = $1, contest_id = $2, participant_type = $3, room = $4, points = $5 
		WHERE user_id = $1
			AND contest_id = $2`,
		participant.UserId, participant.ContestId, participant.ParticipantType, participant.Room, participant.Points)
	if err != nil {
		return err
	}
	return nil
}
