package pgsqlRepository

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	baseRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/base_repo"
)

type SessionnRepository struct {
	baseRepo.BaseRepo
}

func (r *SessionnRepository) List(filter sessionn.Filter) ([]sessionn.Session, error) {
	var sessions []dbSession
	if err := r.DB().Select(&sessions, `
		SELECT *
		FROM sessions
		WHERE ($1 = '' OR $1 = access_token)
	`, filter.AccessToken); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	return toDomainSessions(sessions), nil
}

func (r *SessionnRepository) Upsert(session sessionn.Session) error {
	if _, err := r.DB().NamedExec(`
		INSERT INTO sessions(id, user_id, name, status, access_token, access_expiry, refresh_token, refresh_expiry) 
		VALUES (:id, :user_id, :name, :status, :access_token, :access_expiry, :refresh_token, :refresh_expiry)
	`, toDBSession(session)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	return nil
}

func (r *SessionnRepository) WithTxConn(txConn baseRepo.DbConn) sessionn.Repository {
	return &SessionnRepository{
		BaseRepo: r.BaseRepo.WithTxConn(txConn),
	}
}

func (r *SessionnRepository) InTransaction(fn func(txRepo sessionn.Repository) error) error {
	return baseRepo.InTransaction(r, fn)
}

type dbSession struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	Name          string    `db:"name"`
	Status        string    `db:"status"`
	AccessToken   string    `db:"access_token"`
	AccessExpiry  time.Time `db:"access_expiry"`
	RefreshToken  string    `db:"refresh_token"`
	RefreshExpiry time.Time `db:"refresh_expiry"`
}

func toDBSession(session sessionn.Session) dbSession {
	return dbSession{
		ID:            session.ID.String(),
		UserID:        session.UserID.String(),
		Name:          session.Name,
		Status:        session.Status,
		AccessToken:   session.AccessToken.Token,
		AccessExpiry:  session.AccessToken.Expiry,
		RefreshToken:  session.RefreshToken.Token,
		RefreshExpiry: session.RefreshToken.Expiry,
	}
}

func toDomainSession(session dbSession) sessionn.Session {
	return sessionn.Session{
		ID:     uuid.MustParse(session.ID),
		UserID: uuid.MustParse(session.UserID),
		Name:   session.Name,
		Status: session.Status,
		AccessToken: sessionn.Token{
			Token:  session.AccessToken,
			Expiry: session.AccessExpiry,
		},
		RefreshToken: sessionn.Token{
			Token:  session.RefreshToken,
			Expiry: session.RefreshExpiry,
		},
	}
}

func toDomainSessions(sessions []dbSession) []sessionn.Session {
	domainSessions := make([]sessionn.Session, len(sessions))
	for i, s := range sessions {
		domainSessions[i] = toDomainSession(s)
	}

	return domainSessions
}
