package service

import "errors"

var (
	ErrInvalidLogin        = errors.New("некорректное значение Login")
	ErrInvalidToken        = errors.New("некорректное значение Token")
	ErrInvalidSubjectID    = errors.New("некорректное значение SubjectID")
	ErrInvalidChatID       = errors.New("некорректное значение ChatID")
	ErrInvalidInvitationID = errors.New("некорректное значение InvitationID")
	ErrInvalidUserID       = errors.New("некорректное значение UserID")
	ErrInvalidName         = errors.New("некорректное значение Name")
	ErrInvalidPassword     = errors.New("некорректное значение Password")
	ErrInvalidNick         = errors.New("некорректное значение Nick")
	ErrInvalidUserCode     = errors.New("некорректное значение UserCode")
	ErrInvalidProvider     = errors.New("некорректное значение Provider")
)

var (
	ErrWrongUserCode                = errors.New("неправильный UserCode")
	ErrUserIsNotMember              = errors.New("user не является участником чата")
	ErrSubjectIsNotMember           = errors.New("subject user не является участником чата")
	ErrChatNotExists                = errors.New("чата с таким ID не существует")
	ErrMemberCannotDeleteHimself    = errors.New("участник не может удалить самого себя")
	ErrSubjectUserShouldNotBeChief  = errors.New("пользователь является главным администратором чата")
	ErrSubjectUserIsNotChief        = errors.New("пользователь не является главным администратором чата")
	ErrUnauthorizedChatsView        = errors.New("нельзя просматривать чужой список чатов")
	ErrUnauthorizedInvitationsView  = errors.New("нельзя просматривать чужой список приглашений")
	ErrUserNotExists                = errors.New("пользователя не существует")
	ErrUserIsAlreadyInChat          = errors.New("пользователь уже состоит в чате")
	ErrUserIsAlreadyInvited         = errors.New("пользователь уже приглашен в чат")
	ErrInvitationNotExists          = errors.New("приглашения не существует")
	ErrSubjectUserNotAllowed        = errors.New("у пользователя нет прав на это действие")
	ErrLoginIsAlreadyInUse          = errors.New("логин уже используется")
	ErrProvidersUserIsAlreadyLinked = errors.New("пользователь OAuth-провайдера уже связан с пользователем")
	ErrUnknownOAuthProvider         = errors.New("неизвестный OAuth провайдер")
	ErrLoginOrPasswordDoesNotMatch  = errors.New("не совпадает Login или Password")
)
