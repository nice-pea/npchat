package service

//
//import (
//	"errors"
//
//	"github.com/saime-0/nice-pea-chat/internal/domain"
//)
//
//type Members struct {
//	MembersRepo domain.MembersRepository
//	History     History
//}
//
//type MemberWasDeleted struct {
//	Member domain.Member
//}
//
//type MemberDeleteErr struct {
//	Error error
//}
//
//var (
//	ErrMemberNotFound = errors.New("member not found")
//)
//
//func (m *Members) Delete(id string) error {
//	// Найти участника по ID
//	members, err := m.MembersRepo.List(domain.MembersFilter{ID: id})
//	if err != nil {
//		m.History.Write(MemberDeleteErr{Error: err})
//		return err
//	}
//
//	// Проверить удалось ли найти одного участника
//	if len(members) != 1 {
//		m.History.Write(MemberDeleteErr{Error: ErrMemberNotFound})
//		return ErrMemberNotFound
//	}
//	member := members[0]
//
//	// Удалить участника
//	if err = m.MembersRepo.Delete(id); err != nil {
//		m.History.Write(MemberDeleteErr{Error: err})
//		return err
//	}
//
//	m.History.Write(MemberWasDeleted{Member: member})
//
//	return nil
//}
//
//type MembersListFilter struct {
//	ID      string
//	UserID  string
//	ChatID  string
//	IsOwner *bool
//}
//
//func (m *Members) List(filter MembersListFilter) ([]domain.Member, error) {
//	return m.MembersRepo.List(domain.MembersFilter{
//		ID:      filter.ID,
//		UserID:  filter.UserID,
//		ChatID:  filter.ChatID,
//		IsOwner: filter.IsOwner,
//	})
//}
