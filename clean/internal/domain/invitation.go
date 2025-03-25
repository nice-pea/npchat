package domain

type Invitation struct {
	ID string
	// UserID string
	ChatID string `db:"chat_id"`
}

type InvitationsRepository interface {
	List(filter InvitationsFilter) ([]Invitation, error)
	Save(invitation Invitation) error
	Delete(id string) error
}

type InvitationsFilter struct {
	ID     string
	ChatID string 
}
