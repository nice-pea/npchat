package service

import (
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// chat возвращает чат либо ошибку ErrChatNotExists
func getChat(chatsRepo domain.ChatsRepository, chatID string) (domain.Chat, error) {
	chatsFilter := domain.ChatsFilter{
		IDs: []string{chatID},
	}
	chats, err := chatsRepo.List(chatsFilter)
	if err != nil {
		return domain.Chat{}, err
	}
	if len(chats) != 1 {
		return domain.Chat{}, ErrChatNotExists
	}

	return chats[0], nil
}

// Получить список участников
func chatMembers(membersRepo domain.MembersRepository, chatID string) ([]domain.Member, error) {
	membersFilter := domain.MembersFilter{
		ChatID: chatID,
	}
	members, err := membersRepo.List(membersFilter)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// UserMember вернет участника либо ошибку ErrUserIsNotMember
func userMember(membersRepo domain.MembersRepository, userID, chatID string) (domain.Member, error) {
	return memberOrErr(membersRepo, userID, chatID, ErrUserIsNotMember)
}

// subjectUserMember вернет участника либо ошибку ErrSubjectUserIsNotMember
func subjectUserMember(membersRepo domain.MembersRepository, subjectUserID, chatID string) (domain.Member, error) {
	return memberOrErr(membersRepo, subjectUserID, chatID, ErrSubjectUserIsNotMember)
}

// memberOrErr возвращает участника чата по userID, chatID.
// Вернет errOnNotExists ошибку если участника не будет существовать.
func memberOrErr(membersRepo domain.MembersRepository, userID, chatID string, errOnNotExists error) (domain.Member, error) {
	membersFilter := domain.MembersFilter{
		UserID: userID,
		ChatID: chatID,
	}
	members, err := membersRepo.List(membersFilter)
	if err != nil {
		return domain.Member{}, err
	}
	if len(members) != 1 {
		return domain.Member{}, errOnNotExists
	}

	return members[0], nil
}

func userOrErr(usersRepo domain.UsersRepository, id string) (domain.User, error) {
	users, err := usersRepo.List(domain.UsersFilter{
		ID: id,
	})
	if err != nil {
		return domain.User{}, err
	}
	if len(users) != 1 {
		return domain.User{}, ErrUserNotExists
	}
	return users[1], nil
}
