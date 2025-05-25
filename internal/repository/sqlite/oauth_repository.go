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
	Provider     string    `db:"provider"`
}

type link struct {
	ID         string `db:"id"`
	UserID     string `db:"user_id"`
	ExternalID string `db:"external_id"`
	Provider   string `db:"provider"`
}

func tokenFromDomain(t domain.OAuthToken) token {
	return token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		Expiry:       t.Expiry,
		LinkID:       t.LinkID,
		Provider:     t.Provider,
	}
}

func linkFromDomain(l domain.OAuthLink) link {
	return link{
		UserID:     l.UserID,
		ExternalID: l.ExternalID,
		Provider:   l.Provider,
	}
}

func linksToDomain(links []link) []domain.OAuthLink {
	result := make([]domain.OAuthLink, len(links))
	for i, l := range links {
		result[i] = domain.OAuthLink{
			UserID:     l.UserID,
			ExternalID: l.ExternalID,
			Provider:   l.Provider,
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
		INSERT INTO oauth_tokens (access_token, token_type, refresh_token, expiry, link_id, provider)
		VALUES (:access_token, :token_type, :refresh_token, :expiry, :link_id, :provider)
	`, tokenFromDomain(token))
	if err != nil {
		return fmt.Errorf("DB.NamedExec: %w", err)
	}

	return nil
}

func (r *OAuthRepository) SaveLink(link domain.OAuthLink) error {
	if link == (domain.OAuthLink{}) {
		return errors.New("link  must not be empty")
	}

	_, err := r.DB.NamedExec(`	
		INSERT OR REPLACE INTO oauth_links (user_id, external_id, provider)
		VALUES (:user_id, :external_id, :provider)
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
			WHERE ($1 = '' OR $1 = user_id)
				AND ($2 = '' OR $2 = external_id)
				AND ($3 = '' OR $3 = provider)
		`, filter.UserID, filter.ExternalID, filter.Provider); err != nil {
		return nil, fmt.Errorf("DB.Select: %w", err)
	}

	return linksToDomain(links), nil
}
