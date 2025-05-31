package chatt

type Repository interface {
	ByChatFilter(ChatsFilter) ([]Chat, error)
	ByParticipantsFilter(MembersFilter) ([]Chat, error)
	ByInvitationsFilter(InvitationsFilter) ([]Chat, error)
	Upsert(Chat) error
}

// ChatsFilter представляет собой фильтр по чатам.
type ChatsFilter struct {
	IDs []string // Список идентификаторов чатов для фильтрации
}

// MembersFilter представляет собой фильтр по членам чата.
type MembersFilter struct {
	ID     string // Фильтрация по ID
	UserID string // Фильтрация по пользователю
	ChatID string // Фильтрация по чату
}

// InvitationsFilter представляет собой фильтр по приглашениям.
type InvitationsFilter struct {
	ID            string // Фильтрация по ID приглашения
	ChatID        string // Фильтрация по чату
	UserID        string // Фильтрация по приглашаемому пользователю
	SubjectUserID string // Фильтрация по пригласившему пользователю
}

//type CAFilter struct {
//	IDs                 []string
//	InvolvedUsers       []string
//	HasInvitations      []string // Фильтрация по ID приглашения
//	InvitedUsers        []string // Фильтрация по приглашаемому пользователю
//	SentInvitationUsers []string // Фильтрация по пригласившему пользователю
//}
