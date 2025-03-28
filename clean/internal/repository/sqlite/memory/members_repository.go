package memory

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type member struct {
	ID     string `db:"id"`
	ChatID string `db:"chat_id"`
	UserID string `db:"user_id"`
}

func memberToDomain(repoMember member) domain.Member {
	return domain.Member{
		ID:     repoMember.ID,
		UserID: repoMember.UserID,
		ChatID: repoMember.ChatID,
	}
}

func memberFromDomain(domainMember domain.Member) member {
	return member{
		ID:     domainMember.ID,
		ChatID: domainMember.ChatID,
		UserID: domainMember.UserID,
	}
}

func membersToDomain(repoMembers []member) []domain.Member {
	members := make([]domain.Member, len(repoMembers))
	for i, repoMember := range repoMembers {
		members[i] = memberToDomain(repoMember)
	}
	return members
}

func membersFromDomain(domainMembers []domain.Member) []member {
	repoMembers := make([]member, len(domainMembers))
	for i, repoMember := range domainMembers {
		repoMembers[i] = memberFromDomain(repoMember)
	}
	return repoMembers
}
func (m *SQLiteInMemory) NewMembersRepository() (domain.MembersRepository, error) {
	return &MembersRepository{
		DB: m.db,
	}, nil
}

// Эта строка будет выдавать ошибку во время компиляции,
// если тип справа не имплементирует интерфейс слева
//var _ domain.MembersRepository = (*MembersRepository)(nil) // TODO: раскомментить.

type MembersRepository struct {
	DB *sqlx.DB
}

func (c *MembersRepository) List(filter domain.MembersFilter) ([]domain.Member, error) {
	members := make([]member, 0)
	if err := c.DB.Select(&members, `
			SELECT * 
			FROM members 
			WHERE ($1 = "" OR $1 = id)
				AND ($2 = "" OR $2 = chat_id)
				AND ($3 = "" OR $3 = user_id)
		`, filter.ID, filter.ChatID, filter.UserID); err != nil {
		return nil, fmt.Errorf("error selecting members: %w", err)
	}

	return membersToDomain(members), nil
}

func (c *MembersRepository) Save(member domain.Member) error {
	if member.ID == "" {
		return fmt.Errorf("invalid member id")
	}
	_, err := c.DB.Exec(`
		INSERT INTO members(id, chat_id, user_id)
		VALUES (?, ?, ?)`,
		member.ID, member.ChatID, member.UserID)
	if err != nil {
		return fmt.Errorf("error inserting member: %w", err)
	}

	return nil
}

func (c *MembersRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("invalid member id")
	}
	_, err := c.DB.Exec("DELETE FROM members WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting member: %w", err)
	}

	return nil
}
