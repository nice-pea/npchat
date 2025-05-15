package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// MyInvitations регистрирует обработчик, позволяющий получить список приглашений пользователя.
// Доступен только авторизованным пользователям.
// Метод: GET /invitations
func MyInvitations(router http2.Router) {
	router.HandleFunc(
		"GET /invitations",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			// Формируем входные данные для получения приглашений пользователя.
			// SubjectUserID - ID пользователя, чьи приглашения запрашиваются (сам пользователь).
			// UserID - ID пользователя, выполняющего запрос (сам пользователь).
			input := service.UserInvitationsInput{
				SubjectUserID: context.Session().UserID,
				UserID:        context.Session().UserID,
			}

			// Вызываем сервис получения приглашений пользователя и возвращаем результат.
			return context.Services().Invitations().UserInvitations(input)
		})
}
