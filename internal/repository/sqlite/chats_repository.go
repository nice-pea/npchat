package sqlite

//type chat struct {
//	ID          string `db:"id"`
//	Name        string `db:"name"`
//	ChiefUserID string `db:"chief_user_id"`
//}
//
//func chatToDomain(repoChat chat) domain.Chat {
//	return domain.Chat{
//		ID:          repoChat.ID,
//		Name:        repoChat.Name,
//		ChiefUserID: repoChat.ChiefUserID,
//	}
//}
//
//func chatFromDomain(domainChat domain.Chat) chat {
//	return chat{
//		ID:          domainChat.ID,
//		Name:        domainChat.Name,
//		ChiefUserID: domainChat.ChiefUserID,
//	}
//}
//
//func chatsToDomain(repoChats []chat) []domain.Chat {
//	domainChats := make([]domain.Chat, len(repoChats))
//	for i, repoChat := range repoChats {
//		domainChats[i] = chatToDomain(repoChat)
//	}
//
//	return domainChats
//}
//
//func (r *ChatsRepository) List(filter domain.ChatsFilter) ([]domain.Chat, error) {
//	// Построить запрос, используя bqb
//	where := bqb.Optional("WHERE")
//	if len(filter.IDs) > 0 {
//		where.And("id IN (?)", filter.IDs)
//	}
//	sql, args, err := bqb.New("SELECT * FROM chats ?", where).ToSql()
//	if err != nil {
//		return nil, err
//	}
//
//	// Выполнить запрос, используя sqlx
//	chats := make([]chat, 0)
//	if err = r.DB.Select(&chats, sql, args...); err != nil {
//		return nil, fmt.Errorf("DB.Select: %w", err)
//	}
//
//	return chatsToDomain(chats), nil
//}
//
//func (r *ChatsRepository) Save(chat domain.Chat) error {
//	if chat.ID == "" {
//		return fmt.Errorf("invalid chat id")
//	}
//	_, err := r.DB.NamedExec(`
//		INSERT OR REPLACE INTO chats(id, name, chief_user_id)
//		VALUES (:id, :name, :chief_user_id)
//	`, chatFromDomain(chat))
//	if err != nil {
//		return fmt.Errorf("DB.NamedExec: %w", err)
//	}
//
//	return nil
//}
//
//func (r *ChatsRepository) Delete(id string) error {
//	if id == "" {
//		return fmt.Errorf("invalid chat id")
//	}
//	_, err := r.DB.Exec(`DELETE FROM chats WHERE id = ?`, id)
//	if err != nil {
//		return fmt.Errorf("DB.Exec: %w", err)
//	}
//
//	return nil
//}
