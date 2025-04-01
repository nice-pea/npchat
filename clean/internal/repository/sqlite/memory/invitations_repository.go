package memory

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type invitation struct {
	ID            string `db:"id"`
	UserID        string `db:"user_id"`
	ChatID        string `db:"chat_id"`
	SubjectUserID string `db:"subject_user_id"`
}

func invitationToDomain(repoInvitation invitation) domain.Invitation {
	return domain.Invitation{
		ID:            repoInvitation.ID,
		UserID:        repoInvitation.UserID,
		ChatID:        repoInvitation.ChatID,
		SubjectUserID: repoInvitation.SubjectUserID,
	}
}

func invitationsToDomain(repoInvitations []invitation) []domain.Invitation {
	domainInvitations := make([]domain.Invitation, len(repoInvitations))
	for i, repoInv := range repoInvitations {
		domainInvitations[i] = invitationToDomain(repoInv)
	}
	return domainInvitations
}

func invitationFromDomain(domainInvitation domain.Invitation) invitation {
	return invitation{
		ID:            domainInvitation.ID,
		UserID:        domainInvitation.UserID,
		ChatID:        domainInvitation.ChatID,
		SubjectUserID: domainInvitation.SubjectUserID,
	}
}

func invitationsFromDomain(domainInvitations []domain.Invitation) []invitation {
	repoInvitations := make([]invitation, len(domainInvitations))
	for i, domainInv := range domainInvitations {
		repoInvitations[i] = invitationFromDomain(domainInv)
	}
	return repoInvitations
}

func (m *SQLiteInMemory) NewInvitationsRepository() (domain.InvitationsRepository, error) {
	return &InvitationsRepository{
		DB: m.db,
	}, nil
}

type InvitationsRepository struct {
	DB *sqlx.DB
}

func (c *InvitationsRepository) List(filter domain.InvitationsFilter) ([]domain.Invitation, error) {
	invitations := make([]invitation, 0)
	if err := c.DB.Select(&invitations, `
			SELECT * 
			FROM invitations 
			WHERE ($1 = "" OR $1 = id)
				AND ($2 = "" OR $2 = chat_id)
				AND ($3 = "" OR $3 = user_id)
				AND ($4 = "" OR $4 = subject_user_id)
		`, filter.ID, filter.ChatID, filter.UserID, filter.SubjectUserID); err != nil {
		return nil, fmt.Errorf("error selecting chats: %w", err)
	}

	return invitationsToDomain(invitations), nil
}

func (c *InvitationsRepository) Save(invitation domain.Invitation) error {
	if invitation.ID == "" {
		return fmt.Errorf("invalid invitation id")
	}
	_, err := c.DB.Exec(
		"INSERT INTO invitations(id, chat_id, user_id, subject_user_id) VALUES (?, ?, ?, ?)",
		invitation.ID, invitation.ChatID, invitation.UserID, invitation.SubjectUserID)
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
