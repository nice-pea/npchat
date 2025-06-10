package chatt

import "errors"

var (
	ErrInvalidChiefID                     = errors.New("некорректное значение ChiefID")
	ErrInvalidChatName                    = errors.New("некорректный Name")
	ErrInvalidUserID                      = errors.New("некорректное значение UserID")
	ErrParticipantNotExists               = errors.New("участника не существует")
	ErrSubjectIsNotMember                 = errors.New("subject user не является участником чата")
	ErrCannotRemoveChief                  = errors.New("нельзя удалить главного администратора")
	ErrParticipantExists                  = errors.New("пользователь уже состоит в чате")
	ErrUserIsAlreadyInvited               = errors.New("пользователь уже приглашен в чат")
	ErrInvitationNotExists                = errors.New("приглашения не существует")
	ErrSubjectAndRecipientMustBeDifferent = errors.New("subject и recipient не могут быть одним лицом")
	ErrChatNotExists                      = errors.New("чата с таким ID не существует")
)
