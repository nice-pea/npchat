package chatInvitations

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

var (
	ErrInvalidSubjectID   = errors.New("некорректное значение SubjectID")
	ErrSubjectIsNotMember = errors.New("subject user не является участником чата")
	ErrInvalidChatID      = errors.New("некорректное значение ChatID")
)

// ChatInvitationsIn параметры для запроса приглашений конкретного чата
type In struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
}

// Validate валидирует параметры для запроса приглашений конкретного чата
func (in In) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// ChatInvitationsOut результат запроса приглашений конкретного чата
type Out struct {
	Invitations []chatt.Invitation
}

type ChatInvitationsUsecase struct {
	Repo chatt.Repository
}

// ChatInvitations возвращает список приглашений в конкретный чат.
// Если SubjectID является администратором, то возвращается все приглашения в данный чат,
// иначе только те приглашения, которые отправил именно пользователь.
func (c *ChatInvitationsUsecase) ChatInvitations(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return Out{}, err
	}

	// Проверить является ли пользователь участником чата
	if !chat.HasParticipant(in.SubjectID) {
		return Out{}, ErrSubjectIsNotMember
	}

	// Сохранить сначала все приглашения
	invitations := chat.Invitations

	// Если пользователь не является администратором,
	// то оставить только те приглашения, которые отправил именно пользователь.
	if chat.ChiefID != in.SubjectID {
		invitations = chat.SubjectInvitations(in.SubjectID)
	}

	return Out{
		Invitations: invitations,
	}, err
}
