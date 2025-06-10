package service

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

// Chats сервис, объединяющий случаи использования(юзкейсы) в контексте агрегата чатов
type Chats struct {
	Repo chatt.Repository
}

// WhichParticipateIn входящие параметры
type WhichParticipateIn struct {
	SubjectID string
	UserID    string // TODO: удалить
}

// Validate валидирует значение отдельно каждого параметры
func (in WhichParticipateIn) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// WhichParticipateOut результат запроса чатов
type WhichParticipateOut struct {
	Chats []chatt.Chat
}

// WhichParticipate возвращает список чатов, в которых участвует пользователь
func (c *Chats) WhichParticipate(in WhichParticipateIn) (WhichParticipateOut, error) {
	// Валидировать параметры
	var err error
	if err = in.Validate(); err != nil {
		return WhichParticipateOut{}, err
	}

	// Пользователь может запрашивать только свой список чатов
	if in.UserID != in.SubjectID {
		return WhichParticipateOut{}, ErrUnauthorizedChatsView
	}

	// Получить список участников с фильтром по пользователю
	chats, err := c.Repo.List(chatt.Filter{
		ParticipantID: in.UserID,
	})
	if err != nil {
		return WhichParticipateOut{}, err
	}

	return WhichParticipateOut{Chats: chats}, err
}

// CreateChatIn входящие параметры
type CreateChatIn struct {
	Name        string
	ChiefUserID string // TODO: переименовать в SubjectID
}

func (in CreateChatIn) Validate() error {
	if err := domain.ValidateID(in.ChiefUserID); err != nil {
		return errors.Join(ErrInvalidChiefID, err)
	}
	if err := chatt.ValidateChatName(in.Name); err != nil {
		return errors.Join(ErrInvalidName, err)
	}

	return nil
}

// CreateChatOut результат создания чата
type CreateChatOut struct {
	Chat chatt.Chat
}

// CreateChat создает новый чат и участника для главного администратора - пользователя, который создал этот чат
func (c *Chats) CreateChat(in CreateChatIn) (CreateChatOut, error) {
	if err := in.Validate(); err != nil {
		return CreateChatOut{}, err
	}

	chat, err := chatt.NewChat(in.Name, in.ChiefUserID)
	if err != nil {
		return CreateChatOut{}, err
	}

	// Сохранить чат в репозиторий
	if err := c.Repo.Upsert(chat); err != nil {
		return CreateChatOut{}, err
	}

	return CreateChatOut{
		Chat: chat,
	}, nil
}

// UpdateNameIn входящие параметры
type UpdateNameIn struct {
	SubjectID string
	ChatID    string
	NewName   string
}

// Validate валидирует значение отдельно каждого параметры
func (in UpdateNameIn) Validate() error {
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

// UpdateNameOut результат обновления названия чата
type UpdateNameOut struct {
	Chat chatt.Chat
}

// UpdateName обновляет название чата.
// Доступно только для главного администратора этого чата
func (c *Chats) UpdateName(in UpdateNameIn) (UpdateNameOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return UpdateNameOut{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return UpdateNameOut{}, err
	}

	// Проверить доступ пользователя к этому действию
	if in.SubjectID != chat.ChiefID {
		return UpdateNameOut{}, ErrSubjectUserIsNotChief
	}

	// Перезаписать с новым значением
	if err = chat.UpdateName(in.NewName); err != nil {
		return UpdateNameOut{}, err
	}
	if err = c.Repo.Upsert(chat); err != nil {
		return UpdateNameOut{}, err
	}

	return UpdateNameOut{
		Chat: chat,
	}, nil
}
