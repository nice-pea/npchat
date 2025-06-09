package sqlite

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type session struct {
	ID     string `db:"id"`
	UserID string `db:"user_id"`
	Token  string `db:"token"`
	Status int    `db:"status"`
}

func sessionFromDomain(s domain.Session) session {
	return session{
		ID:     s.ID,
		UserID: s.UserID,
		Token:  s.Token,
		Status: s.Status,
	}
}

func sessionToDomain(s session) domain.Session {
	return domain.Session{
		ID:     s.ID,
		UserID: s.UserID,
		Token:  s.Token,
		Status: s.Status,
	}
}

func sessionsToDomain(sessions []session) []domain.Session {
	result := make([]domain.Session, len(sessions))
	for i, s := range sessions {
		result[i] = sessionToDomain(s)
	}
	return result
}

//func (r *SessionsRepository) List(filter domain.SessionsFilter) ([]domain.Session, error) {
//	sessions := make([]session, 0)
//	if err := r.DB.Select(&sessions, `
//			SELECT *
//			FROM sessions
//			WHERE ($1 = '' OR $1 = token)
//		`, filter.Token); err != nil {
//		return nil, fmt.Errorf("DB.Select: %w", err)
//	}
//
//	return sessionsToDomain(sessions), nil
//}
//
//func (r *SessionsRepository) Save(session domain.Session) error {
//	if session.ID == "" {
//		return errors.New("invalid session id")
//	}
//	_, err := r.DB.NamedExec(`
//		INSERT OR REPLACE INTO sessions (id, user_id, token, status)
//		VALUES (:id, :user_id, :token, :status)
//	`, sessionFromDomain(session))
//	if err != nil {
//		return fmt.Errorf("DB.NamedExec: %w", err)
//	}
//
//	return nil
//}
//
//func (r *SessionsRepository) Delete(id string) error {
//	if id == "" {
//		return errors.New("invalid session id")
//	}
//	_, err := r.DB.Exec(`DELETE FROM sessions WHERE id = ?`, id)
//	if err != nil {
//		return fmt.Errorf("DB.Exec: %w", err)
//	}
//
//	return nil
//}
