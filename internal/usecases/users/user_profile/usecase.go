package userProfile

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

var (
	ErrInvalidSubjectID = errors.New("некорректное значение SubjectID")
	ErrInvalidUserID    = errors.New("некорректное значение UserID")
)

// In входящие параметры
type In struct {
	SubjectID uuid.UUID
	UserID    uuid.UUID
}

// validate валидирует значение отдельно каждого параметра
func (in In) validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return ErrInvalidUserID
	}

	return nil
}

// Out результат запроса чатов
type Out struct {
	User userr.User
}

type UserProfileUsecase struct {
	Repo userr.Repository
}

// UserProfile возвращает информацию о пользователе
func (c *UserProfileUsecase) UserProfile(in In) (Out, error) {
	// Валидировать параметры
	if err := in.validate(); err != nil {
		return Out{}, err
	}

	// Получить пользователя
	user, err := userr.Find(c.Repo, userr.Filter{
		ID: in.UserID,
	})
	if err != nil {
		return Out{}, err
	}

	// Очистить чувствительные данные, если запрашивается чужой профиль
	if user.ID != in.SubjectID {
		user.OpenAuthUsers = nil
		user.BasicAuth = userr.BasicAuth{}
	}

	return Out{
		User: user,
	}, nil
}
