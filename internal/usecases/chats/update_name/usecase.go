package updateName

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
	"github.com/nice-pea/npchat/internal/usecases/events"
)

var (
	ErrInvalidSubjectID      = errors.New("некорректное значение SubjectID")
	ErrInvalidChatID         = errors.New("некорректное значение ChatID")
	ErrInvalidName           = errors.New("некорректное значение Name")
	ErrSubjectUserIsNotChief = errors.New("пользователь не является главным администратором чата")
)

// In входящие параметры
type In struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
	NewName   string
}

// Validate валидирует значение отдельно каждого параметры
func (in In) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := chatt.ValidateChatName(in.NewName); err != nil {
		return errors.Join(err, ErrInvalidName)
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// Out результат обновления названия чата
type Out struct {
	Chat chatt.Chat
}

type UpdateNameUsecase struct {
	Repo          chatt.Repository
	EventConsumer events.Consumer
}

// UpdateName обновляет название чата.
// Доступно только для главного администратора этого чата
func (c *UpdateNameUsecase) UpdateName(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return Out{}, err
	}

	// Проверить доступ пользователя к этому действию
	if in.SubjectID != chat.ChiefID {
		return Out{}, ErrSubjectUserIsNotChief
	}

	// Инициализировать буфер событий
	eventsBuf := new(events.Buffer)

	// Перезаписать с новым значением
	if err = chat.UpdateName(in.NewName, eventsBuf); err != nil {
		return Out{}, err
	}
	if err = c.Repo.Upsert(chat); err != nil {
		return Out{}, err
	}

	// Отправить собранные события
	c.EventConsumer.Consume(eventsBuf.Events())

	return Out{
		Chat: chat,
	}, nil
}
