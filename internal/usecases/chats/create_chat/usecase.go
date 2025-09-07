package createChat

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
	"github.com/nice-pea/npchat/internal/usecases/events"
)

var (
	ErrInvalidChiefID = errors.New("некорректное значение ChiefID")
	ErrInvalidName    = errors.New("некорректное значение Name")
)

// In входящие параметры
type In struct {
	Name        string
	ChiefUserID uuid.UUID // TODO: переименовать в SubjectID
}

func (in In) Validate() error {
	if err := domain.ValidateID(in.ChiefUserID); err != nil {
		return errors.Join(ErrInvalidChiefID, err)
	}
	if err := chatt.ValidateChatName(in.Name); err != nil {
		return errors.Join(ErrInvalidName, err)
	}

	return nil
}

// Out результат создания чата
type Out struct {
	Chat chatt.Chat
}

type CreateChatUsecase struct {
	Repo          chatt.Repository
	EventConsumer events.Consumer
}

// CreateChat создает новый чат и участника для главного администратора - пользователя, который создал этот чат
func (c *CreateChatUsecase) CreateChat(in In) (Out, error) {
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Инициализировать буфер событий
	eventsBuf := new(events.Buffer)

	chat, err := chatt.NewChat(in.Name, in.ChiefUserID, eventsBuf)
	if err != nil {
		return Out{}, err
	}

	// Сохранить чат в репозиторий
	if err := c.Repo.Upsert(chat); err != nil {
		return Out{}, err
	}

	// Отправить собранные события
	c.EventConsumer.Consume(eventsBuf.Events())

	return Out{
		Chat: chat,
	}, nil
}
