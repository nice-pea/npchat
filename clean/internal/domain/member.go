package domain

type Member struct {
	ID string
	//UserID string
	ChatID string
}

func (m Member) ValidateID() error {
	panic("unimplemented")
}

type MembersRepository interface {
	List(filter MembersFilter) ([]Member, error)
	Save(member Member) error
	Delete(id string) error
}

type MembersFilter struct {
	ID string
	//UserID string
	ChatID string
	//IsOwner *bool
}
