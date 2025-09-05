package completeOauthLogin

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

// In представляет собой параметры завершения входа через Oauth.
type In struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера Oauth
}

// Out представляет собой результат завершения входа через Oauth.
type Out struct {
	Session sessionn.Session // Сессия пользователя
	User    userr.User       // Пользователь
}

type CompleteOauthLoginUsecase struct {
}

// CompleteOauthLogin завершает процесс входа пользователя через Oauth.
func (u *CompleteOauthLoginUsecase) CompleteOauthLogin(in In) (Out, error) {
	// TODO
	return Out{}, nil
}
