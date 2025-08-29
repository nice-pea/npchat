package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

// Chats сервис, объединяющий случаи использования(юзкейсы) в контексте агрегата чатов
type Chats struct {
	Repo chatt.Repository
}

// WhichParticipateIn входящие параметры
type WhichParticipateIn struct {
	SubjectID uuid.UUID
	UserID    uuid.UUID // TODO: удалить
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