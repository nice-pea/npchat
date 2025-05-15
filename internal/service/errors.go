package service

import "errors"

var (
	ErrInvalidLogin                = errors.New("некорректный Login")
	ErrInvalidToken                = errors.New("некорректный Token")
	ErrInvalidID                   = errors.New("некорректный ID")
	ErrInvalidSubjectUserID        = errors.New("некорректный SubjectUserID")
	ErrInvalidChatID               = errors.New("некорректный ChatID")
	ErrInvalidInvitationID         = errors.New("некорректный InvitationID")
	ErrInvalidUserID               = errors.New("некорректный UserID")
	ErrInvalidName                 = errors.New("некорректный Name")
	ErrInvalidChiefUserID          = errors.New("некорректный ChiefUserID")
	ErrUserIsNotMember             = errors.New("user не является участником чата")
	ErrSubjectUserIsNotMember      = errors.New("subject user не является участником чата")
	ErrChatNotExists               = errors.New("чата с таким ID не существует")
	ErrMemberCannotDeleteHimself   = errors.New("участник не может удалить самого себя")
	ErrSubjectUserShouldNotBeChief = errors.New("пользователь является главным администратором чата")
	ErrSubjectUserIsNotChief       = errors.New("пользователь не является главным администратором чата")
	ErrUnauthorizedChatsView       = errors.New("нельзя просматривать чужой список чатов")
	ErrUnauthorizedInvitationsView = errors.New("нельзя просматривать чужой список приглашений")
	ErrUserNotExists               = errors.New("пользователя не существует")
	ErrUserIsAlreadyInChat         = errors.New("пользователь уже состоит в чате")
	ErrUserIsAlreadyInvited        = errors.New("пользователь уже приглашен в чат")
	ErrInvitationNotExists         = errors.New("приглашения не существует")
	ErrSubjectUserNotAllowed       = errors.New("у пользователя нет прав на это действие")
)
