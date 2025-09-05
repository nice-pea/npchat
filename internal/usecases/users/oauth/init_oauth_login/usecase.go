package initOauthLogin

import (
	"errors"
)

var (
	ErrInvalidSubjectID      = errors.New("некорректное значение SubjectID")
	ErrInvalidName           = errors.New("некорректное значение Name")
	ErrInvalidUserID         = errors.New("некорректное значение UserID")
	ErrUnauthorizedChatsView = errors.New("нельзя просматривать чужой список чатов")
)

// In представляет собой параметры инициализации входа через Oauth.
type In struct {
	Provider string // Имя провайдера Oauth
}

// Out представляет собой результат инициализации входа через Oauth.
type Out struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
}

type InitOauthLoginUsecase struct {
}

// InitOauthLogin инициализирует процесс входа пользователя через Oauth.
func (u *InitOauthLoginUsecase) InitOauthLogin(in In) (Out, error) {
	// TODO
	return Out{}, nil
}
