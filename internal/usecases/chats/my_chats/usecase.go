package myChats

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvalidSubjectID      = errors.New("некорректное значение SubjectID")
	ErrInvalidName           = errors.New("некорректное значение Name")
	ErrInvalidUserID         = errors.New("некорректное значение UserID")
	ErrUnauthorizedChatsView = errors.New("нельзя просматривать чужой список чатов")
)

// In входящие параметры
type In struct {
	SubjectID uuid.UUID
	UserID    uuid.UUID // TODO: удалить
}

// Validate валидирует значение отдельно каждого параметры
func (in In) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// Out результат запроса чатов
type Out struct {
	Chats []chatt.Chat
}

type MyChatsUsecase struct {
	Repo chatt.Repository
}

// MyChats возвращает список чатов, в которых участвует пользователь
func (c *MyChatsUsecase) MyChats(in In) (Out, error) {
	// Валидировать параметры
	var err error
	if err = in.Validate(); err != nil {
		return Out{}, err
	}

	// Пользователь может запрашивать только свой список чатов
	if in.UserID != in.SubjectID {
		return Out{}, ErrUnauthorizedChatsView
	}

	// Получить список участников с фильтром по пользователю
	chats, err := c.Repo.List(chatt.Filter{
		ParticipantID: in.UserID,
	})
	if err != nil {
		return Out{}, err
	}

	return Out{Chats: chats}, err
}
