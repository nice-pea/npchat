package chatMembers

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvalidSubjectID   = errors.New("некорректное значение SubjectID")
	ErrInvalidChatID      = errors.New("некорректное значение ChatID")
	ErrSubjectIsNotMember = errors.New("subject user не является участником чата")
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
type Out struct {
	Participants []chatt.Participant
}

type ChatMembersUsecase struct {
	Repo chatt.Repository
}

// ChatMembers возвращает список участников чата
func (c *ChatMembersUsecase) ChatMembers(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return Out{}, err
	}

	// Пользователь должен быть участником чата
	if !chat.HasParticipant(in.SubjectID) {
		return Out{}, ErrSubjectIsNotMember
	}

	return Out{
		Participants: chat.Participants,
	}, nil
}
