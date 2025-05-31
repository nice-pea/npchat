package domain

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type ChatAggregate struct {
	ID      string // Уникальный ID чата
	Name    string // Название чата
	ChiefID string // ID главного пользователя чата

	Participants []Participant
	Invitations  []Invitation2
}

type Participant struct {
	UserID string // ID пользователя, который является участником чата
}

type Invitation2 struct {
	ID          string // Глобальный уникальный ID приглашения
	RecipientID string // Пользователь, получивший приглашение
	SubjectID   string // Пользователь, отправивший приглашение
}

func NewChat(name string, chiefID string) (ChatAggregate, error) {
	if err := ValidateChatName(name); err != nil {
		return ChatAggregate{}, err
	}
	if err := ValidateChiefID(chiefID); err != nil {
		return ChatAggregate{}, errors.Join(err, ErrInvalidChiefUserID)
	}

	return ChatAggregate{
		ID:      uuid.NewString(),
		Name:    name,
		ChiefID: chiefID,
		Participants: []Participant{
			{UserID: chiefID}, // Главный администратор
		},
		Invitations: nil,
	}, nil
}

type ChatAggregateRepository interface {
	ByChatFilter(ChatsFilter) ([]ChatAggregate, error)
	ByParticipantsFilter(MembersFilter) ([]ChatAggregate, error)
	ByInvitationsFilter(InvitationsFilter) ([]ChatAggregate, error)
	Upsert(ChatAggregate) error
}

//type CAFilter struct {
//	IDs                 []string
//	InvolvedUsers       []string
//	HasInvitations      []string // Фильтрация по ID приглашения
//	InvitedUsers        []string // Фильтрация по приглашаемому пользователю
//	SentInvitationUsers []string // Фильтрация по пригласившему пользователю
//}

func NewParticipant(userID string) (Participant, error) {
	if err := ValidateID(userID); err != nil {
		return Participant{}, err
	}

	return Participant{
		UserID: userID,
	}, nil
}

func (c *ChatAggregate) HasParticipant(userID string) bool {
	for _, p := range c.Participants {
		if p.UserID == userID {
			return true
		}
	}

	return false
}

func (c *ChatAggregate) HasInvitationWithRecipient(recipientID string) bool {
	for _, i := range c.Invitations {
		if i.RecipientID == recipientID {
			return true
		}
	}

	return false
}

func (c *ChatAggregate) HasInvitation(id string) bool {
	for _, i := range c.Invitations {
		if i.ID == id {
			return true
		}
	}

	return false
}

//func (c *ChatAggregate) RemoveInvitationByRecipient(recipientID string) error {
//	if !c.HasInvitationWithRecipient(recipientID) {
//		return ErrInvitationNotExists
//	}
//
//	c.Invitations = slices.DeleteFunc(c.Invitations, func(i Invitation2) bool {
//		return i.RecipientID == recipientID
//	})
//
//	return nil
//}

func (c *ChatAggregate) RemoveInvitation(id string) error {
	if !c.HasInvitation(id) {
		return ErrInvitationNotExists
	}

	c.Invitations = slices.DeleteFunc(c.Invitations, func(i Invitation2) bool {
		return i.ID == id
	})

	return nil
}

func (c *ChatAggregate) Invitation(id string) (Invitation2, error) {
	for _, i := range c.Invitations {
		if i.ID == id {
			return i, nil
		}
	}

	return Invitation2{}, ErrInvitationNotExists
}

func (c *ChatAggregate) RemoveParticipant(userID string) error {
	if userID == c.ChiefID {
		return ErrSubjectUserShouldNotBeChief
	}

	if !c.HasParticipant(userID) {
		return ErrUserIsNotMember
	}

	c.Participants = slices.DeleteFunc(c.Participants, func(p Participant) bool {
		return p.UserID == userID
	})

	return nil
}

func (c *ChatAggregate) AddParticipant(p Participant) error {
	// Проверить является ли subject участником чата
	if c.HasParticipant(p.UserID) {
		return ErrUserIsAlreadyInChat
	}

	// Проверить, не существует ли приглашение для этого пользователя в этот чат
	if c.HasInvitationWithRecipient(p.UserID) {
		return ErrUserIsAlreadyInvited
	}

	c.Participants = append(c.Participants, p)

	return nil
}

func NewInvitation(subjectID, recipientID string) (Invitation2, error) {
	if err := ValidateID(subjectID); err != nil {
		return Invitation2{}, err
	}
	if err := ValidateID(recipientID); err != nil {
		return Invitation2{}, err
	}
	if recipientID == subjectID {
		// Subject и User не могут быть одним лицом
		return Invitation2{}, ErrSubjectAndRecipientMustBeDifferent
	}

	return Invitation2{
		ID:          uuid.NewString(),
		RecipientID: recipientID,
		SubjectID:   subjectID,
	}, nil
}

// AddInvitation добавляет приглашение в чат
func (c *ChatAggregate) AddInvitation(invitation Invitation2) error {
	// Проверить является ли subject участником чата
	if !c.HasParticipant(invitation.SubjectID) {
		return ErrSubjectIsNotMember
	}

	// Проверить является ли user участником чата
	if c.HasParticipant(invitation.RecipientID) {
		return ErrUserIsAlreadyInChat
	}

	// Проверить, не существует ли приглашение для этого пользователя в этот чат
	if c.HasInvitationWithRecipient(invitation.RecipientID) {
		return ErrUserIsAlreadyInvited
	}

	c.Invitations = append(c.Invitations, invitation)

	return nil
}

func (c *ChatAggregate) UpdateName(name string) error {
	if err := ValidateChatName(name); err != nil {
		return err
	}

	c.Name = name

	return nil
}
