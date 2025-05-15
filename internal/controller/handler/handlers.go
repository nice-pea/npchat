package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Получить список чатов
func RegisterMyChatsHandler(router http2.Router) {
	router.HandleFunc(
		"GET /chats",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.UserChatsInput{
				SubjectUserID: context.Session().UserID,
				UserID:        context.Session().UserID,
			}

			chats, err := context.Services().Chats().UserChats(input)
			if err != nil {
				return nil, err
			}

			return chats, nil
		})
}

// Создать чат
func RegisterCreateChatHandler(router http2.Router) {
	type requestBody struct {
		Name string `json:"name"`
		//ChiefUserID string `json:"chief_user_id"`
	}
	router.HandleFunc(
		"POST /chats",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.CreateInput{
				ChiefUserID: context.Session().UserID,
				Name:        rb.Name,
			}

			result, err := context.Services().Chats().Create(input)
			if err != nil {
				return nil, err
			}

			return result, nil
		})
}

// Обновить название чата
func RegisterUpdateChatNameHandler(router http2.Router) {
	router.HandleFunc(
		"PUT /chats/{chatID}/name",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Покинуть чат
func RegisterLeaveChatHandler(router http2.Router) {
	router.HandleFunc(
		"POST /chats/{chatID}/leave",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Получить список участников чата
func RegisterChatMembersHandler(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/members",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Удалить участника из чата
func RegisterDeleteMemberHandler(router http2.Router) {
	router.HandleFunc(
		"DELETE /chats/{chatID}/members/{memberID}",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Получить список моих приглашений
func RegisterMyInvitationsHandler(router http2.Router) {
	router.HandleFunc(
		"GET /invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Получить список приглашений в чат
func RegisterChatInvitationsHandler(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Отправить приглашение в чат
func RegisterSendInvitationHandler(router http2.Router) {
	router.HandleFunc(
		"POST /invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Принять приглашение в чат
func RegisterAcceptInvitationHandler(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/accept",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}

// Отменить приглашение в чат
func RegisterCancelInvitationHandler(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/cancel",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
