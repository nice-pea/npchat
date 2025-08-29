package cancelInvitation

import (
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvitationNotExists   = errors.New("приглашения не существует")
	ErrInvalidSubjectID      = errors.New("некорректное значение SubjectID")
	ErrInvalidInvitationID   = errors.New("некорректное значение InvitationID")
	ErrSubjectIsNotMember    = errors.New("subject user не является участником чата")
	ErrSubjectUserNotAllowed = errors.New("у пользователя нет прав на это действие")
)

type In struct {
	SubjectID    uuid.UUID
	InvitationID uuid.UUID
}

func (in In) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}

	return nil
}

type Out struct{}

type CancelInvitationUsecase struct {
	Repo chatt.Repository
}

// CancelInvitation отменяет приглашение
func (c *CancelInvitationUsecase) CancelInvitation(in In) (Out, error) {
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

	// Достать приглашение из чата
	invitation, err := chat.Invitation(in.InvitationID)
	if err != nil {
		return Out{}, err
	}

	if in.SubjectID == invitation.SubjectID {
		// Проверить, существование участника чата
		if !chat.HasParticipant(invitation.SubjectID) {
			return Out{}, ErrSubjectIsNotMember
		}
	}

	// Список тех, кто может отменять приглашение
	allowedSubjects := []uuid.UUID{
		chat.ChiefID,           // Главный администратор
		invitation.SubjectID,   // Пригласивший
		invitation.RecipientID, // Приглашаемый
	}
	// Проверить, может ли пользователь отменить приглашение
	if !slices.Contains(allowedSubjects, in.SubjectID) {
		return Out{}, ErrSubjectUserNotAllowed
	}

	// Удаляем приглашение из чата
	if err := chat.RemoveInvitation(in.InvitationID); err != nil {
		return Out{}, err
	}

	// Сохранить чат в репозиторий
	if err := c.Repo.Upsert(chat); err != nil {
		return Out{}, err
	}

	return Out{}, nil
}
