package sqlite

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type token struct {
	AccessToken  string    `db:"access_token"`
	TokenType    string    `db:"token_type"`
	RefreshToken string    `db:"refresh_token"`
	Expiry       time.Time `db:"expiry"`
	LinkID       string    `db:"link_id"`
}

type link struct {
	ID         string `db:"id"`
	UserID     string `db:"user_id"`
	ExternalID string `db:"external_id"`
}

func tokenFromDomain(t domain.OAuthToken) token {
	return token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		LinkID:       t.LinkID,
	}
}

func linkFromDomain(l domain.OAuthLink) link {
	return link{
		ID:         l.ID,
		UserID:     l.UserID,
		ExternalID: l.ExternalID,
	}
}

func linksToDomain(links []link) []domain.OAuthLink {
	result := make([]domain.OAuthLink, len(links))
	for i, l := range links {
		result[i] = domain.OAuthLink{
			ID:         l.ID,
			UserID:     l.UserID,
			ExternalID: l.ExternalID,
		}
	}
	return result
}

type OAuthRepository struct {
	DB *sqlx.DB
}

func (m *RepositoryFactory) NewOAuthRepository() domain.OAuthRepository {
	return &OAuthRepository{
		DB: m.db,
	}
}

func (r *OAuthRepository) SaveToken(token domain.OAuthToken) error {
	if token == (domain.OAuthToken{}) {
		return errors.New("token must not be empty")
	}

	_, err := r.DB.NamedExec(`	
		INSERT INTO oauth_tokens (access_token, token_type, refresh_token, expiry, link_id)
		VALUES (:access_token, :token_type, :refresh_token, :expiry, :link_id)
	`, tokenFromDomain(token))
	if err != nil {
		return fmt.Errorf("DB.NamedExec: %w", err)
	}

	return nil
}

func (r *OAuthRepository) SaveLink(link domain.OAuthLink) error {
	if link.ID == "" {
		return errors.New("invalid id")
	}

	_, err := r.DB.NamedExec(`	
		INSERT INTO oauth_links (id, user_id, external_id)
		VALUES (:id, :user_id, :external_id)
	`, linkFromDomain(link))
	if err != nil {
		return fmt.Errorf("DB.NamedExec: %w", err)
	}

	return err
}

func (r *OAuthRepository) ListLinks(filter domain.OAuthListLinksFilter) ([]domain.OAuthLink, error) {
	links := make([]link, 0)
	if err := r.DB.Select(&links, `
			SELECT * 
			FROM oauth_links 
			WHERE ($1 = '' OR $1 = id)
				AND ($2 = '' OR $2 = user_id)
				AND ($3 = '' OR $3 = external_id)
		`, filter.ID, filter.UserID, filter.ExternalID); err != nil {
		return nil, fmt.Errorf("DB.Select: %w", err)
	}

	return linksToDomain(links), nil
}
