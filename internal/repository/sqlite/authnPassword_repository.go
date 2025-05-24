package sqlite

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (m *RepositoryFactory) NewAuthnPasswordRepository() domain.AuthnPasswordRepository {
	return &AuthnPasswordRepository{
		DB: m.db,
	}
}

type AuthnPasswordRepository struct {
	DB *sqlx.DB
}

type authnPassword struct {
	UserID   string `db:"user_id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func authnPasswordFromDomain(c domain.AuthnPassword) authnPassword {
	return authnPassword{
		UserID:   c.UserID,
		Login:    c.Login,
		Password: c.Password,
	}
}

func authnPasswordToDomain(c authnPassword) domain.AuthnPassword {
	return domain.AuthnPassword{
		UserID:   c.UserID,
		Login:    c.Login,
		Password: c.Password,
	}
}

func authnPasswordListToDomain(authnPasswords []authnPassword) []domain.AuthnPassword {
	result := make([]domain.AuthnPassword, len(authnPasswords))
	for i, s := range authnPasswords {
		result[i] = authnPasswordToDomain(s)
	}
	return result
}

func (r *AuthnPasswordRepository) List(filter domain.AuthnPasswordFilter) ([]domain.AuthnPassword, error) {
	aps := make([]authnPassword, 0)
	if err := r.DB.Select(&aps, `
			SELECT * 
			FROM authn_passwords 
			WHERE ($1 = '' OR $1 = user_id)
				AND ($2 = '' OR $2 = login)
				AND ($3 = '' OR $3 = password)
		`, filter.UserID, filter.Login, filter.Password); err != nil {
		return nil, fmt.Errorf("DB.Select: %w", err)
	}

	return authnPasswordListToDomain(aps), nil
}

func (r *AuthnPasswordRepository) Save(ap domain.AuthnPassword) error {
	if ap.UserID == "" {
		return errors.New("invalid user id")
	}
	_, err := r.DB.NamedExec(`
		INSERT OR REPLACE INTO authn_passwords (user_id, login, password)
		VALUES (:user_id, :login, :password)
	`, authnPasswordFromDomain(ap))
	if err != nil {
		return fmt.Errorf("DB.NamedExec: %w", err)
	}
	return nil
}

func (r *AuthnPasswordRepository) GetByLogin(login string) (domain.AuthnPassword, error) {
	if login == "" {
		return domain.AuthnPassword{}, errors.New("login is empty")
	}
	var ap authnPassword
	err := r.DB.Get(&ap, `
		SELECT * FROM authn_passwords WHERE login = ?
	`, login)
	if err != nil {
		return domain.AuthnPassword{}, fmt.Errorf("DB.Get: %w", err)
	}
	return authnPasswordToDomain(ap), nil
}

func (r *AuthnPasswordRepository) DeleteByUserID(userID string) error {
	if userID == "" {
		return errors.New("userID is empty")
	}
	_, err := r.DB.Exec(`DELETE FROM authn_passwords WHERE user_id = ?`, userID)
	if err != nil {
		return fmt.Errorf("DB.Exec: %w", err)
	}
	return nil
}
