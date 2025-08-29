package leaveChat

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvalidSubjectID = errors.New("некорректное значение SubjectID")
	ErrInvalidChatID    = errors.New("некорректное значение ChatID")
)

// In входящие параметры
type In struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
}

// Validate валидирует значение отдельно каждого параметры
func (in In) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}

	return nil
}

// Out результат запроса чатов
type Out struct{}

type LeaveChatUsecase struct {
	Repo chatt.Repository
}

// LeaveChat удаляет участника из чата
func (c *LeaveChatUsecase) LeaveChat(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return Out{}, err
	}

	// Удалить пользователя из чата
	if err = chat.RemoveParticipant(in.SubjectID); err != nil {
		return Out{}, err
	}

	// Сохранить чат в репозиторий
	if err = c.Repo.Upsert(chat); err != nil {
		return Out{}, err
	}

	return Out{}, nil
}
