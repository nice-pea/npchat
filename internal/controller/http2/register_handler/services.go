package register_handler

import (
	"github.com/nice-pea/npchat/internal/service"
)

// services определяет интерфейс для доступа к сервисам приложения
type services interface {
	Chats() *service.Chats       // Сервис чатов
	Sessions() *service.Sessions // Сервис сессий
	Users() *service.Users       // Сервис пользователей
}
