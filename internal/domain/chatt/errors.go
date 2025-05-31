package chatt

import "errors"

var (
	ErrInvalidChiefID                     = errors.New("некорректное значение ChiefID")
	ErrInvalidChatName                    = errors.New("некорректный Name")
	ErrUserIsNotMember                    = errors.New("user не является участником чата")
	ErrSubjectIsNotMember                 = errors.New("subject user не является участником чата")
	ErrSubjectUserShouldNotBeChief        = errors.New("пользователь является главным администратором чата")
	ErrUserIsAlreadyInChat                = errors.New("пользователь уже состоит в чате")
	ErrUserIsAlreadyInvited               = errors.New("пользователь уже приглашен в чат")
	ErrInvitationNotExists                = errors.New("приглашения не существует")
	ErrSubjectAndRecipientMustBeDifferent = errors.New("subject и recipient не могут быть одним лицом")
	ErrChatNotExists                      = errors.New("чата с таким ID не существует")
)
