package register_handler

import (
	"github.com/nice-pea/npchat/internal/service"
)

// Services определяет интерфейс для доступа к сервисам приложения
type Services interface {
	Chats() *service.Chats       // Сервис чатов
	Sessions() *service.Sessions // Сервис сессий
	Users() *service.Users       // Сервис пользователей
}
