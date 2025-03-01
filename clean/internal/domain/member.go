package domain

type Member struct {
	ID     uint
	UserID uint
	ChatID uint
}

type MembersRepository interface {
	List(filter MembersFilter) ([]Member, error)
	Save(member Member) error
	Delete(id uint) error
}

type MembersFilter struct {
	ID     uint
	UserID uint
	ChatID uint
}
