package pgsqlRepository

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	sqlxRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/sqlx_repo"
)

type ChattRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *ChattRepository) List(filter chatt.Filter) ([]chatt.Chat, error) {
	// Запросить чаты
	var chats []dbChat
	if err := r.DB().Select(&chats, `
		SELECT c.*
		FROM chats c
			LEFT JOIN participants p
			    ON c.id = p.chat_id
			LEFT JOIN invitations i
				ON c.id = i.chat_id
		WHERE ($1 = '' OR $1 = c.id)
			AND ($2 = '' OR $2 = i.id)
			AND ($3 = '' OR $3 = i.recipient_id)
			AND ($4 = '' OR $4 = p.user_id)
		GROUP BY c.id
	`, filter.ID, filter.InvitationID, filter.InvitationRecipientID,
		filter.ParticipantID); err != nil {
		return nil, err
	}

	// Если чатов нет, сразу вернуть пустой список
	if len(chats) == 0 {
		return nil, nil
	}

	// Собрать ID найденных чатов
	chatIDs := make([]string, len(chats))
	for i, c := range chats {
		chatIDs[i] = c.ID
	}

	// Найти участников чатов
	var participants []dbParticipant
	if err := r.DB().Select(&participants, `
		SELECT *
		FROM participants
		WHERE chat_id = ANY($1)
	`, chatIDs); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	// Создать карту, где ключ это ID чата, а значение это список его участников
	participantsMap := make(map[string][]dbParticipant, len(chats))
	for _, p := range participants {
		participantsMap[p.ChatID] = append(participantsMap[p.ChatID], p)
	}

	// Найти приглашения в чат
	var invitations []dbInvitation
	if err := r.DB().Select(&invitations, `
		SELECT *
		FROM invitations
		WHERE chat_id = ANY($1)
	`, chatIDs); err != nil {
		return nil, fmt.Errorf("r.DB().Select: %w", err)
	}

	// Создать карту, где ключ это ID чата, а значение это список приглашений в него
	invitationsMap := make(map[string][]dbInvitation, len(chats))
	for _, i := range invitations {
		invitationsMap[i.ChatID] = append(invitationsMap[i.ChatID], i)
	}

	return toDomainChats(chats, participantsMap, invitationsMap), nil
}

func (r *ChattRepository) Upsert(chat chatt.Chat) error {
	if chat.ID == uuid.Nil {
		return fmt.Errorf("chat ID is required")
	}

	if r.IsTx() {
		return r.upsert(chat)
	} else {
		return r.InTransaction(func(txRepo chatt.Repository) error {
			return txRepo.Upsert(chat)
		})
	}
}

func (r *ChattRepository) upsert(chat chatt.Chat) error {
	if _, err := r.DB().NamedExec(`
		INSERT INTO chats(id, name, chief_id) 
		VALUES (:id, :name, :chief_id)
		ON CONFLICT DO UPDATE SET
			name=excluded.name,
			chief_id=excluded.chief_id
	`, toDBChat(chat)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	if _, err := r.DB().NamedExec(`
		DELETE
		FROM participants
		WHERE chat_id = :chat_id;
		INSERT INTO participants(chat_id, user_id)
		VALUES (:chat_id, :user_id)
	`, toDBParticipants(chat)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	if _, err := r.DB().NamedExec(`
		DELETE
		FROM invitations
		WHERE chat_id = :chat_id;
		INSERT INTO invitations(id, chat_id, subject_id, recipient_id)
		VALUES (:id, :chat_id, :subject_id, :recipient_id)
	`, toDBInvitations(chat)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	return nil
}

func (r *ChattRepository) InTransaction(fn func(txRepo chatt.Repository) error) error {
	return r.SqlxRepo.InTransaction(func(txSqlxRepo sqlxRepo.SqlxRepo) error {
		return fn(&ChattRepository{SqlxRepo: txSqlxRepo})
	})
}

type dbChat struct {
	ID      string `db:"id"`
	Name    string `db:"name"`
	ChiefID string `db:"chief_id"`
}

func toDBChat(chat chatt.Chat) dbChat {
	return dbChat{
		ID:      chat.ID.String(),
		Name:    chat.Name,
		ChiefID: chat.ChiefID.String(),
	}
}

func toDomainChat(
	chat dbChat,
	participants []dbParticipant,
	invitations []dbInvitation,
) chatt.Chat {
	return chatt.Chat{
		ID:           uuid.MustParse(chat.ID),
		Name:         chat.Name,
		ChiefID:      uuid.MustParse(chat.ChiefID),
		Participants: toDomainParticipants(participants),
		Invitations:  toDomainInvitations(invitations),
	}
}

func toDomainChats(
	chats []dbChat,
	participants map[string][]dbParticipant,
	invitations map[string][]dbInvitation,
) []chatt.Chat {
	domainChats := make([]chatt.Chat, len(chats))
	for i, chat := range chats {
		domainChats[i] = toDomainChat(chat, participants[chat.ID], invitations[chat.ID])
	}

	return domainChats
}

type dbParticipant struct {
	ChatID string `db:"chat_id"`
	UserID string `db:"user_id"`
}

func toDBParticipants(chat chatt.Chat) []dbParticipant {
	dbParticipants := make([]dbParticipant, len(chat.Participants))
	for i, p := range chat.Participants {
		dbParticipants[i] = dbParticipant{
			ChatID: chat.ID.String(),
			UserID: p.UserID.String(),
		}
	}

	return dbParticipants
}
func toDomainParticipants(participants []dbParticipant) []chatt.Participant {
	pp := make([]chatt.Participant, len(participants))
	for i, p := range participants {
		pp[i] = chatt.Participant{
			UserID: uuid.MustParse(p.UserID),
		}
	}

	return pp
}

type dbInvitation struct {
	ID          string `db:"id"`
	ChatID      string `db:"chat_id"`
	SubjectID   string `db:"subject_id"`
	RecipientID string `db:"recipient_id"`
}

func toDBInvitations(chat chatt.Chat) []dbInvitation {
	dbInvitations := make([]dbInvitation, len(chat.Invitations))
	for i, inv := range chat.Invitations {
		dbInvitations[i] = dbInvitation{
			ID:          inv.ID.String(),
			ChatID:      chat.ID.String(),
			SubjectID:   inv.SubjectID.String(),
			RecipientID: inv.RecipientID.String(),
		}
	}

	return dbInvitations
}

func toDomainInvitations(invitations []dbInvitation) []chatt.Invitation {
	ii := make([]chatt.Invitation, len(invitations))
	for i, inv := range invitations {
		ii[i] = chatt.Invitation{
			ID:          uuid.MustParse(inv.ID),
			RecipientID: uuid.MustParse(inv.RecipientID),
			SubjectID:   uuid.MustParse(inv.SubjectID),
		}
	}

	return ii
}
