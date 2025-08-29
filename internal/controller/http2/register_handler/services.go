package register_handler

import (
	"github.com/nice-pea/npchat/internal/service"
	sessionsFind "github.com/nice-pea/npchat/internal/service/sessions_find"
)

// Services определяет интерфейс для доступа к сервисам приложения
type Services interface {
	Chats() *service.Chats // Сервис чато
	SessionsFind(sessionsFind.In) (sessionsFind.Out, error)
	Users() *service.Users // Сервис пользователей
}
