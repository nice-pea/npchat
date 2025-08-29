package receivedInvitations

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvalidSubjectID = errors.New("некорректное значение SubjectID")
)

type In struct {
	SubjectID uuid.UUID
}

func (in In) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// Out входящие параметры
type Out struct {
	// ChatsInvitations карта приглашений, где ключ - chatID, значение - приглашение
	ChatsInvitations map[uuid.UUID]chatt.Invitation
}

type ReceivedInvitationsUsecase struct {
	Repo chatt.Repository
}

// ReceivedInvitations возвращает список приглашений конкретного пользователя в чаты
func (c *ReceivedInvitationsUsecase) ReceivedInvitations(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти чат
	chats, err := c.Repo.List(chatt.Filter{
		InvitationRecipientID: in.SubjectID,
	})
	if err != nil {
		return Out{}, err
	}

	// Если нет чатов, вернут пустой список
	if len(chats) == 0 {
		return Out{}, nil
	}

	// Собрать приглашения, полученные пользователем
	invitations := make(map[uuid.UUID]chatt.Invitation, len(chats))
	for _, chat := range chats {
		invitations[chat.ID], _ = chat.RecipientInvitation(in.SubjectID)
	}

	return Out{
		ChatsInvitations: invitations,
	}, err
}
