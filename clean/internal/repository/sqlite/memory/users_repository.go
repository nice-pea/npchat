package memory

import (
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

func (m *SQLiteInMemory) NewUsersRepository() (domain.UsersRepository, error) {
	return &UsersRepository{
		DB: m.db,
	}, nil
}

type UsersRepository struct {
	DB *sqlx.DB
}

func (c *UsersRepository) List(filter domain.UsersFilter) ([]domain.User, error) {
	users := make([]user, 0)
	if err := c.DB.Select(&users, `
			SELECT * 
			FROM users 
			WHERE ($1 = "" OR $1 = id)
		`, filter.ID); err != nil {
		return nil, fmt.Errorf("error selecting users: %w", err)
	}

	return usersToDomain(users), nil
}

func (c *UsersRepository) Save(user domain.User) error {
	if user.ID == "" {
		return fmt.Errorf("invalid user id")
	}
	_, err := c.DB.Exec(`
		INSERT OR REPLACE INTO users(id)
		VALUES (?)`,
		user.ID)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}

	return nil
}

func (c *UsersRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("invalid user id")
	}
	_, err := c.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
