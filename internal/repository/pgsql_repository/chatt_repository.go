package pgsqlRepository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nullism/bqb"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	sqlxRepo "github.com/nice-pea/npchat/internal/repository/pgsql_repository/sqlx_repo"
)

type ChattRepository struct {
	sqlxRepo.SqlxRepo
}

func (r *ChattRepository) List(filter chatt.Filter) ([]chatt.Chat, error) {
	sel := bqb.New("SELECT c.* FROM chats c")
	where := bqb.Optional("WHERE")

	needJoinParticipants := filter.ParticipantID != uuid.Nil
	if needJoinParticipants {
		sel = sel.Space("LEFT JOIN participants p ON c.id = p.chat_id")
	}
	if filter.ParticipantID != uuid.Nil {
		where = where.And("p.user_id = ?", filter.ParticipantID)
	}

	needJoinInvitations := filter.InvitationID != uuid.Nil || filter.InvitationRecipientID != uuid.Nil
	if needJoinInvitations {
		sel = sel.Space("LEFT JOIN invitations i ON c.id = i.chat_id")
	}
	if filter.InvitationID != uuid.Nil {
		where = where.And("i.id = ?", filter.InvitationID)
	}
	if filter.InvitationRecipientID != uuid.Nil {
		where = where.And("i.recipient_id = ?", filter.InvitationRecipientID)
	}

	if filter.ID != uuid.Nil {
		where = where.And("c.id = ?", filter.ID)
	}

	if !filter.ActiveBefore.IsZero() {
		where = where.And("c.last_active_at < ?", filter.ActiveBefore)
	}

	limit := bqb.New("")
	if filter.Limit > 0 {
		limit = limit.Space("LIMIT ?", filter.Limit)
	}

	query, args, err := bqb.New("? ? GROUP BY c.id ORDER BY last_active_at DESC ?", sel, where, limit).ToPgsql()
	if err != nil {
		return nil, fmt.Errorf("bqb.ToPgsql: %w", err)
	}

	// Запросить чаты
	var chats []dbChat
	if err := r.DB().Select(&chats, query, args...); err != nil {
		return nil, fmt.Errorf("bqb.Select: %w", err)
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
	`, pq.Array(chatIDs)); err != nil {
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
	`, pq.Array(chatIDs)); err != nil {
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
		INSERT INTO chats(id, name, chief_id, last_active_at) 
		VALUES (:id, :name, :chief_id, :last_active_at)
		ON CONFLICT (id) DO UPDATE SET
			name=excluded.name,
			chief_id=excluded.chief_id
	`, toDBChat(chat)); err != nil {
		return fmt.Errorf("r.DB().NamedExec: %w", err)
	}

	// Удалить прошлых участников
	if _, err := r.DB().Exec(`
		DELETE FROM participants WHERE chat_id = $1
	`, chat.ID); err != nil {
		return fmt.Errorf("r.DB().Exec: %w", err)
	}

	if len(chat.Participants) > 0 {
		if _, err := r.DB().NamedExec(`
			INSERT INTO participants(chat_id, user_id)
			VALUES (:chat_id, :user_id)
		`, toDBParticipants(chat)); err != nil {
			return fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	}

	// Удалить прошлые приглашения
	if _, err := r.DB().Exec(`
		DELETE FROM invitations WHERE chat_id = $1
	`, chat.ID); err != nil {
		return fmt.Errorf("r.DB().Exec: %w", err)
	}

	if len(chat.Invitations) > 0 {
		if _, err := r.DB().NamedExec(`
		INSERT INTO invitations(id, chat_id, subject_id, recipient_id)
		VALUES (:id, :chat_id, :subject_id, :recipient_id)
	`, toDBInvitations(chat)); err != nil {
			return fmt.Errorf("r.DB().NamedExec: %w", err)
		}
	}

	return nil
}

func (r *ChattRepository) InTransaction(fn func(txRepo chatt.Repository) error) error {
	return r.SqlxRepo.InTransaction(func(txSqlxRepo sqlxRepo.SqlxRepo) error {
		return fn(&ChattRepository{SqlxRepo: txSqlxRepo})
	})
}

type dbChat struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	ChiefID      string    `db:"chief_id"`
	LastActiveAt time.Time `db:"last_active_at"`
}

func toDBChat(chat chatt.Chat) dbChat {
	return dbChat{
		ID:           chat.ID.String(),
		Name:         chat.Name,
		ChiefID:      chat.ChiefID.String(),
		LastActiveAt: chat.LastActiveAt,
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
		LastActiveAt: chat.LastActiveAt.UTC(),
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
