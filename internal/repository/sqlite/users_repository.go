package sqlite

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type user struct {
	ID string `db:"id"`
}

func userToDomain(repoUser user) domain.User {
	return domain.User{
		ID: repoUser.ID,
	}
}

func userFromDomain(domainUser domain.User) user {
	return user{
		ID: domainUser.ID,
	}
}

func usersToDomain(repoUsers []user) []domain.User {
	users := make([]domain.User, len(repoUsers))
	for i, repoUser := range repoUsers {
		users[i] = userToDomain(repoUser)
	}
	return users
}

func usersFromDomain(domainUsers []domain.User) []user {
	repoUsers := make([]user, len(domainUsers))
	for i, repoUser := range domainUsers {
		repoUsers[i] = userFromDomain(repoUser)
	}
	return repoUsers
}

func (m *RepositoryFactory) NewUsersRepository() domain.UsersRepository {
	return &UsersRepository{
		DB: m.db,
	}
}

type UsersRepository struct {
	DB *sqlx.DB
}

func (r *UsersRepository) List(filter domain.UsersFilter) ([]domain.User, error) {
	users := make([]user, 0)
	if err := r.DB.Select(&users, `
			SELECT * 
			FROM users 
			WHERE ($1 = '' OR $1 = id)
		`, filter.ID); err != nil {
		return nil, fmt.Errorf("DB.Select: %w", err)
	}

	return usersToDomain(users), nil
}

func (r *UsersRepository) Save(user domain.User) error {
	if user.ID == "" {
		return errors.New("invalid user id")
	}
	_, err := r.DB.NamedExec(`
		INSERT OR REPLACE INTO users(id)
		VALUES (:id)
	`, user)
	if err != nil {
		return fmt.Errorf("DB.NamedExec: %w", err)
	}

	return nil
}

func (r *UsersRepository) Delete(id string) error {
	if id == "" {
		return errors.New("invalid user id")
	}
	_, err := r.DB.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("DB.Exec: %w", err)
	}

	return nil
}
