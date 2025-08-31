package acceptInvitation

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvitationNotExists = errors.New("приглашения не существует")
	ErrInvalidSubjectID    = errors.New("некорректное значение SubjectID")
	ErrInvalidInvitationID = errors.New("некорректное значение InvitationID")
)

type In struct {
	SubjectID    uuid.UUID
	InvitationID uuid.UUID
}

func (in In) Validate() error {
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}

	return nil
}

type Out struct{}

type AcceptInvitationUsecase struct {
	Repo chatt.Repository
}

// AcceptInvitation добавляет пользователя в чат, путем принятия приглашения
func (c *AcceptInvitationUsecase) AcceptInvitation(in In) (Out, error) {
	// Валидировать входные данные
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{
		InvitationID: in.InvitationID,
	})
	if errors.Is(err, chatt.ErrChatNotExists) {
		return Out{}, ErrInvitationNotExists
	} else if err != nil {
		return Out{}, err
	}

	// Удаляем приглашение из чата
	if err := chat.RemoveInvitation(in.InvitationID); err != nil {
		return Out{}, err
	}

	// Создаем участника чата
	participant, err := chatt.NewParticipant(in.SubjectID)
	if err != nil {
		return Out{}, err
	}

	// Добавить участника в чат
	if err := chat.AddParticipant(participant); err != nil {
		return Out{}, err
	}

	// Сохранить чат в репозиторий
	if err := c.Repo.Upsert(chat); err != nil {
		return Out{}, err
	}

	return Out{}, nil
}
