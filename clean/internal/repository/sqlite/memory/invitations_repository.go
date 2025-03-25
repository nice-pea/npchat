package memory

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (m *SQLiteInMemory) NewInvitationsRepository() (domain.InvitationsRepository, error) {
	return &InvitationsRepository{
		DB: m.db,
	}, nil
}

type InvitationsRepository struct {
	DB *sqlx.DB
}

func (c *InvitationsRepository) List(filter domain.InvitationsFilter) ([]domain.Invitation, error) {
	invitations := make([]domain.Invitation, 0)
	if err := c.DB.Select(&invitations, `
			SELECT * 
			FROM invitations 
			WHERE ($1 = "" OR $1 = id)
				AND ($2 = "" OR $2 = chat_id)
		`, filter.ID, filter.ChatID); err != nil {
		return nil, fmt.Errorf("error selecting chats: %w", err)
	}

	return invitations, nil
}

func (c *InvitationsRepository) Save(invitation domain.Invitation) error {
	if invitation.ID == "" {
		return fmt.Errorf("invalid invitation id")
	}
	_, err := c.DB.Exec("INSERT INTO invitations(id, chat_id) VALUES (?, ?)", invitation.ID, invitation.ChatID)
	if err != nil {
		return fmt.Errorf("error inserting invitation: %w", err)
	}

	return nil
}

func (c *InvitationsRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("invalid invitation id")
	}
	_, err := c.DB.Exec("DELETE FROM invitations WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting invitation: %w", err)
	}

	return nil
}
