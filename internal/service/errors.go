package service

import "errors"

var (
	ErrInvalidLogin        = errors.New("некорректное значение BasicAuthLogin")
	ErrInvalidToken        = errors.New("некорректное значение Token")
	ErrInvalidSubjectID    = errors.New("некорректное значение SubjectID")
	ErrInvalidChatID       = errors.New("некорректное значение ChatID")
	ErrInvalidInvitationID = errors.New("некорректное значение InvitationID")
	ErrInvalidUserID       = errors.New("некорректное значение UserID")
	ErrInvalidName         = errors.New("некорректное значение Name")
	ErrInvalidPassword     = errors.New("некорректное значение Password")
	ErrInvalidUserCode     = errors.New("некорректное значение UserCode")
	ErrInvalidProvider     = errors.New("некорректное значение Provider")
	ErrInvalidChiefID      = errors.New("некорректное значение ChiefID")
)

var (
	ErrWrongUserCode                = errors.New("неправильный UserCode")
	ErrSubjectIsNotMember           = errors.New("subject user не является участником чата")
	ErrMemberCannotDeleteHimself    = errors.New("участник не может удалить самого себя")
	ErrSubjectUserIsNotChief        = errors.New("пользователь не является главным администратором чата")
	ErrUnauthorizedChatsView        = errors.New("нельзя просматривать чужой список чатов")
	ErrInvitationNotExists          = errors.New("приглашения не существует")
	ErrSubjectUserNotAllowed        = errors.New("у пользователя нет прав на это действие")
	ErrLoginIsAlreadyInUse          = errors.New("логин уже используется")
	ErrProvidersUserIsAlreadyLinked = errors.New("пользователь OAuth-провайдера уже связан с пользователем")
	ErrUnknownOAuthProvider         = errors.New("неизвестный OAuth провайдер")
	ErrLoginOrPasswordDoesNotMatch  = errors.New("не совпадает BasicAuthLogin или Password")
	ErrLoginIsRequired              = errors.New("login это обязательный параметр")
	ErrPasswordIsRequired           = errors.New("password это обязательный параметр")
	ErrNameIsRequired               = errors.New("name это обязательный параметр")
)
