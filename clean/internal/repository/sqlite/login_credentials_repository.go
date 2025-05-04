package sqlite

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (m *RepositoryFactory) NewLoginCredentialsRepository() domain.LoginCredentialsRepository {
	return &LoginCredentialsRepository{
		DB: m.db,
	}
}

type LoginCredentialsRepository struct {
	DB *sqlx.DB
}

type loginCredentials struct {
	UserID   string `db:"user_id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func loginCredentialsFromDomain(c domain.LoginCredentials) loginCredentials {
	return loginCredentials{
		UserID:   c.UserID,
		Login:    c.Login,
		Password: c.Password,
	}
}

func loginCredentialsToDomain(c loginCredentials) domain.LoginCredentials {
	return domain.LoginCredentials{
		UserID:   c.UserID,
		Login:    c.Login,
		Password: c.Password,
	}
}

func loginCredsListToDomain(credentials []loginCredentials) []domain.LoginCredentials {
	result := make([]domain.LoginCredentials, len(credentials))
	for i, s := range credentials {
		result[i] = loginCredentialsToDomain(s)
	}
	return result
}

func (r *LoginCredentialsRepository) List(filter domain.LoginCredentialsFilter) ([]domain.LoginCredentials, error) {
	credentials := make([]loginCredentials, 0)
	if err := r.DB.Select(&credentials, `
			SELECT * 
			FROM login_credentials 
			WHERE ($1 = '' OR $1 = user_id)
				AND ($2 = '' OR $2 = login)
				AND ($3 = '' OR $3 = password)
		`, filter.UserID, filter.Login, filter.Password); err != nil {
		return nil, fmt.Errorf("DB.Select: %w", err)
	}

	return loginCredsListToDomain(credentials), nil
}

func (r *LoginCredentialsRepository) Save(credentials domain.LoginCredentials) error {
	if credentials.UserID == "" {
		return errors.New("invalid user id")
	}
	_, err := r.DB.NamedExec(`
		INSERT OR REPLACE INTO login_credentials (user_id, login, password)
		VALUES (:user_id, :login, :password)
	`, loginCredentialsFromDomain(credentials))
	if err != nil {
		return fmt.Errorf("DB.NamedExec: %w", err)
	}
	return nil
}

func (r *LoginCredentialsRepository) GetByLogin(login string) (domain.LoginCredentials, error) {
	if login == "" {
		return domain.LoginCredentials{}, errors.New("login is empty")
	}
	var creds loginCredentials
	err := r.DB.Get(&creds, `
		SELECT * FROM login_credentials WHERE login = ?
	`, login)
	if err != nil {
		return domain.LoginCredentials{}, fmt.Errorf("DB.Get: %w", err)
	}
	return loginCredentialsToDomain(creds), nil
}

func (r *LoginCredentialsRepository) DeleteByUserID(userID string) error {
	if userID == "" {
		return errors.New("userID is empty")
	}
	_, err := r.DB.Exec(`DELETE FROM login_credentials WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("DB.Exec: %w", err)
	}
	return nil
}
