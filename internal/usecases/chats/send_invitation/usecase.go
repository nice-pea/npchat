package sendInvitation

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvalidSubjectID = errors.New("некорректное значение SubjectID")
	ErrInvalidChatID    = errors.New("некорректное значение ChatID")
	ErrInvalidUserID    = errors.New("некорректное значение UserID")
)

type In struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
}

func (in In) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return ErrInvalidChatID
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return ErrInvalidUserID
	}

	return nil
}

type Out struct {
	Invitation chatt.Invitation
}

type SendInvitationUsecase struct {
	Repo chatt.Repository
}

// SendInvitation отправляет приглашения пользователю от участника чата
func (c *SendInvitationUsecase) SendInvitation(in In) (Out, error) {
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return Out{}, err
	}

	// Создать приглашение
	inv, err := chatt.NewInvitation(in.SubjectID, in.UserID)
	if err != nil {
		return Out{}, err
	}

	// Добавить приглашение в чат
	if err = chat.AddInvitation(inv); err != nil {
		return Out{}, err
	}

	// Сохранить чат в репозиторий
	if err = c.Repo.Upsert(chat); err != nil {
		return Out{}, err
	}

	return Out{
		Invitation: inv,
	}, nil
}
