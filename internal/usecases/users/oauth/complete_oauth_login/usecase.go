package completeOAuthLogin

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

var (
	ErrInvalidSubjectID      = errors.New("некорректное значение SubjectID")
	ErrInvalidName           = errors.New("некорректное значение Name")
	ErrInvalidUserID         = errors.New("некорректное значение UserID")
	ErrUnauthorizedChatsView = errors.New("нельзя просматривать чужой список чатов")
)

// In представляет собой параметры завершения входа через OAuth.
type In struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера OAuth
}

// Out представляет собой результат завершения входа через OAuth.
type Out struct {
	Session sessionn.Session // Сессия пользователя
	User    userr.User       // Пользователь
}

type CompleteOAuthLoginUsecase struct {
}

// CompleteOAuthLogin завершает процесс входа пользователя через OAuth.
func (u *CompleteOAuthLoginUsecase) CompleteOAuthLogin(in In) (Out, error) {
	// TODO
	return Out{}, nil
}
