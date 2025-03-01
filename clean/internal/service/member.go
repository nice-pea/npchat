package service

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type Members struct {
	MembersRepo domain.MembersRepository
	Hi          History
}

type MemberWasDeleted struct {
	Member domain.Member
}

type MemberDeleteErr struct {
	Error error
}

var (
	ErrMemberNotFound = errors.New("member not found")
)

func (m *Members) Delete(id uint) error {
	// Найти участника по ID
	members, err := m.MembersRepo.List(domain.MembersFilter{ID: id})
	if err != nil {
		m.Hi.Write(MemberDeleteErr{Error: err})
		return err
	}

	// Проверить удалось ли найти одного участника
	if len(members) != 1 {
		m.Hi.Write(MemberDeleteErr{Error: ErrMemberNotFound})
		return ErrMemberNotFound
	}
	member := members[0]

	// Удалить участника
	if err = m.MembersRepo.Delete(id); err != nil {
		m.Hi.Write(MemberDeleteErr{Error: err})
		return err
	}

	m.Hi.Write(MemberWasDeleted{Member: member})

	return nil
}
