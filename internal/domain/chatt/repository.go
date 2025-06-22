package chatt

import "github.com/google/uuid"

// Repository представляет собой интерфейс для работы с репозиторием чатов.
type Repository interface {
	List(Filter) ([]Chat, error)
	Upsert(Chat) error
	InTransaction(func(txRepo Repository) error) error
}

// Filter представляет собой фильтр для выборки чатов.
type Filter struct {
	ID                    uuid.UUID // Фильтрация по ID чата
	InvitationID          uuid.UUID // Фильтрация по ID приглашений в чате
	InvitationRecipientID uuid.UUID // Фильтрация по ID получателей приглашения в чат
	ParticipantID         uuid.UUID // Фильтрация по ID участников в чате
}

// Find возвращает чат либо ошибку ErrChatNotExists
func Find(repo Repository, filter Filter) (Chat, error) {
	chats, err := repo.List(filter)
	if err != nil {
		return Chat{}, err
	}
	if len(chats) != 1 {
		return Chat{}, ErrChatNotExists
	}

	return chats[0], nil
}
