package service

import (
	"errors"
	"slices"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// ChatInvitationsInput параметры для запроса приглашений конкретного чата
type ChatInvitationsInput struct {
	SubjectID string
	ChatID    string
}

// Validate валидирует параметры для запроса приглашений конкретного чата
func (in ChatInvitationsInput) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// ChatInvitationsOutput результат запроса приглашений конкретного чата
type ChatInvitationsOutput struct {
	Invitations []domain.Invitation2
}

// ChatInvitations возвращает список приглашений в конкретный чат.
// Если SubjectID является администратором, то возвращается все приглашения в данный чат,
// иначе только те приглашения, которые отправил именно пользователь.
func (c *Chats) ChatInvitations(in ChatInvitationsInput) (ChatInvitationsOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ChatInvitationsOutput{}, err
	}

	// Проверить существование чата
	chat, err := getChatAggregate(c.ChatAggregateRepo, in.ChatID)
	if err != nil {
		return ChatInvitationsOutput{}, err
	}

	// Проверить является ли пользователь участником чата
	if !chat.HasParticipant(in.SubjectID) {
		return ChatInvitationsOutput{}, ErrSubjectIsNotMember
	}

	var invitations []domain.Invitation2

	// Если пользователь не является администратором,
	// то оставить только те приглашения, которые отправил именно пользователь.
	if chat.ChiefID != in.SubjectID {
		invitations = slices.DeleteFunc(chat.Invitations, func(i domain.Invitation2) bool {
			return i.RecipientID != in.SubjectID
		})
	}

	return ChatInvitationsOutput{
		Invitations: invitations,
	}, err
}

type ReceivedInvitationsInput struct {
	SubjectUserID string
}

func (in ReceivedInvitationsInput) Validate() error {
	if err := domain.ValidateID(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// ReceivedInvitationsOutput входящие параметры
type ReceivedInvitationsOutput struct {
	Invitations []domain.Invitation
}

// ReceivedInvitations возвращает список приглашений конкретного пользователя в чаты
func (c *Chats) ReceivedInvitations(in ReceivedInvitationsInput) (ReceivedInvitationsOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ReceivedInvitationsOutput{}, err
	}

	// Проверить существование чата
	invitationsFilter := domain.InvitationsFilter{
		UserID: in.SubjectUserID,
	}
	chats, err := c.ChatAggregateRepo.ByInvitationsFilter(invitationsFilter)
	if err != nil {
		return ReceivedInvitationsOutput{}, err
	}

	// Если нет чатов, вернут пустой список
	if len(chats) == 0 {
		return ReceivedInvitationsOutput{}, nil
	}

	// Получить список полученных пользователем приглашений
	var invitations []domain.Invitation
	for _, chat := range chats {
		for _, invitation := range chat.Invitations {
			if invitation.RecipientID == in.SubjectUserID {
				invitations = append(invitations, domain.Invitation{
					SubjectUserID: invitation.SubjectID,
					UserID:        invitation.RecipientID,
					ChatID:        chat.ID,
				})
			}
		}
	}

	return ReceivedInvitationsOutput{
		Invitations: invitations,
	}, err
}

type SendInvitationInput struct {
	SubjectID string
	ChatID    string
	UserID    string
}

func (in SendInvitationInput) Validate() error {
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

type SendInvitationOutput struct {
	Invitation domain.Invitation
}

// SendInvitation отправляет приглашения пользователю от участника чата
func (c *Chats) SendInvitation(in SendInvitationInput) (SendInvitationOutput, error) {
	if err := in.Validate(); err != nil {
		return SendInvitationOutput{}, err
	}

	// Проверить существование чата
	chat, err := getChatAggregate(c.ChatAggregateRepo, in.ChatID)
	if err != nil {
		return SendInvitationOutput{}, err
	}

	// Создать приглашение
	inv, err := domain.NewInvitation(in.SubjectID, in.UserID)
	if err != nil {
		return SendInvitationOutput{}, err
	}

	// Добавить приглашение в чат
	if err = chat.AddInvitation(inv); err != nil {
		return SendInvitationOutput{}, err
	}

	// Сохранить чат в репозиторий
	if err = c.ChatAggregateRepo.Upsert(chat); err != nil {
		return SendInvitationOutput{}, err
	}

	return SendInvitationOutput{
		Invitation: domain.Invitation{
			SubjectUserID: inv.SubjectID,
			UserID:        inv.RecipientID,
			ChatID:        chat.ID,
		},
	}, nil
}

type AcceptInvitationInput struct {
	SubjectID    string
	InvitationID string
}

func (in AcceptInvitationInput) Validate() error {
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}

	return nil
}

// AcceptInvitation добавляет пользователя в чат, путем принятия приглашения
func (c *Chats) AcceptInvitation(in AcceptInvitationInput) error {
	// Валидировать входные данные
	if err := in.Validate(); err != nil {
		return err
	}

	// Проверить существование чата
	chat, err := getChatAggregateByInvitationID(c.ChatAggregateRepo, in.InvitationID)
	if err != nil {
		return err
	}

	// Удаляем приглашение из чата
	if err := chat.RemoveInvitation(in.InvitationID); err != nil {
		return err
	}

	// Создаем участника чата
	participant, err := domain.NewParticipant(in.SubjectID)
	if err != nil {
		return err
	}

	// Добавить участника в чат
	if err := chat.AddParticipant(participant); err != nil {
		return err
	}

	// Сохранить чат в репозиторий
	if err := c.ChatAggregateRepo.Upsert(chat); err != nil {
		return err
	}

	return nil
}

type CancelInvitationInput struct {
	SubjectUserID string
	InvitationID  string
}

func (in CancelInvitationInput) Validate() error {
	if err := domain.ValidateID(in.SubjectUserID); err != nil {
		return ErrInvalidSubjectID
	}
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}

	return nil
}

// CancelInvitation отменяет приглашение
func (c *Chats) CancelInvitation(in CancelInvitationInput) error {
	// Валидировать входные данные
	if err := in.Validate(); err != nil {
		return err
	}

	// Проверить существование чата
	chat, err := getChatAggregateByInvitationID(c.ChatAggregateRepo, in.InvitationID)
	if err != nil {
		return err
	}

	// Достать приглашение из чата
	invitation, err := chat.Invitation(in.InvitationID)
	if err != nil {
		return err
	}

	if in.SubjectUserID == invitation.SubjectID {
		// Проверить, существование участника чата
		if !chat.HasParticipant(invitation.SubjectID) {
			return ErrSubjectIsNotMember
		}
	}

	// Список тех, кто может отменять приглашение
	allowedSubjects := []string{
		chat.ChiefID,           // Главный администратор
		invitation.SubjectID,   // Пригласивший
		invitation.RecipientID, // Приглашаемый
	}
	// Проверить, может ли пользователь отменить приглашение
	if !slices.Contains(allowedSubjects, in.SubjectUserID) {
		return ErrSubjectUserNotAllowed
	}

	// Удаляем приглашение из чата
	if err := chat.RemoveInvitation(in.InvitationID); err != nil {
		return err
	}

	// Сохранить чат в репозиторий
	if err := c.ChatAggregateRepo.Upsert(chat); err != nil {
		return err
	}

	return nil
}
