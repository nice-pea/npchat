package base

import "github.com/saime-0/nice-pea-chat/internal/domain"

var _ domain.MembersRepository = (*MembersRepository)(nil)

type MembersRepository struct {
}

func (m *MembersRepository) List(filter domain.MembersFilter) ([]domain.Member, error) {
	return nil, nil
}

func (m *MembersRepository) Save(member domain.Member) error {
	return nil
}

func (m *MembersRepository) Delete(id uint) error {
	return nil
}
