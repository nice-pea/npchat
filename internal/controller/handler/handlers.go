package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

func CreateChat(router http2.Router) {
	type requestBody struct {
		Name        string `json:"name"`
		ChiefUserID string `json:"chief_user_id"`
	}
	router.HandleFunc("POST /chats", func(context http2.Context) (any, error) {
		var rb requestBody
		if err := http2.DecodeBody(context, &rb); err != nil {
			return nil, err
		}

		input := service.CreateInput{
			Name:        rb.Name,
			ChiefUserID: rb.ChiefUserID,
		}

		result, err := context.Services().Chats().Create(input)
		if err != nil {
			return nil, err
		}

		return result, nil
	})
}

func MyChats(router http2.Router) {
	router.HandleFunc("GET /chats", func(context http2.Context) (any, error) {
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

func UpdateChatName(router http2.Router) {
	router.HandleFunc("PUT /chats/{chatID}/name", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func LeaveChat(router http2.Router) {
	router.HandleFunc("POST /chats/{chatID}/leave", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func ChatMembers(router http2.Router) {
	router.HandleFunc("GET /chats/{chatID}/members", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func DeleteMember(router http2.Router) {
	router.HandleFunc("DELETE /chats/{chatID}/members/{memberID}", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func MyInvitations(router http2.Router) {
	router.HandleFunc("GET /invitations", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func ChatInvitations(router http2.Router) {
	router.HandleFunc("GET /chats/{chatID}/invitations", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func SendInvitation(router http2.Router) {
	router.HandleFunc("POST /invitations", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func AcceptInvitation(router http2.Router) {
	router.HandleFunc("POST /invitations/{invitationID}/accept", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func CancelInvitation(router http2.Router) {
	router.HandleFunc("POST /invitations/{invitationID}/cancel", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}
