package initOAuthLogin

import (
	"errors"
)

var (
	ErrInvalidSubjectID      = errors.New("некорректное значение SubjectID")
	ErrInvalidName           = errors.New("некорректное значение Name")
	ErrInvalidUserID         = errors.New("некорректное значение UserID")
	ErrUnauthorizedChatsView = errors.New("нельзя просматривать чужой список чатов")
)

// In представляет собой параметры инициализации входа через OAuth.
type In struct {
	Provider string // Имя провайдера OAuth
}

// Out представляет собой результат инициализации входа через OAuth.
type Out struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
}

type InitOAuthLoginUsecase struct {
}

// InitOAuthLogin инициализирует процесс входа пользователя через OAuth.
func (u *InitOAuthLoginUsecase) InitOAuthLogin(in In) (Out, error) {
	// TODO
	return Out{}, nil
}
