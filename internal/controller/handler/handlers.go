package handler

import (
	"encoding/json"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

func Ping(router http2.Router) {
	router.HandleFunc("/ping", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

// Публичные эндпоинты, без аутентификации

type requestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func LoginByPassword(router http2.Router) {
	router.HandleFunc("POST /login/password", func(context http2.Context) (any, error) {
		var rb requestBody
		if err := http2.DecodeBody(context, &rb); err != nil {
			return nil, err
		}

		input := service.AuthnPasswordLoginInput{
			Login:    rb.Login,
			Password: rb.Password,
		}

		session, err := context.Services().AuthnPassword().Login(input)
		if err != nil {
			return nil, err
		}

		return session, nil
	})
}

// Эндпоинты, требующие аутентификации

func MyChats(router http2.Router) {
	router.HandleFunc("GET /chats", func(context http2.Context) (any, error) {
		input := service.UserChatsInput{
			SubjectUserID: context.session.UserID,
			UserID:        context.session.UserID,
		}

		chats, err := c.chats.UserChats(input)
		if err != nil {
			return nil, err
		}

		return chats, nil
	})
}

func CreateChat(router http2.Router) {
	router.HandleFunc("POST /chats", func(context http2.Context) (any, error) {
		return "pong", nil
	})
}

func (c *Controller) registerHandlers() {
	// Чат
	//c.HandleFunc("GET /chats", c.MyChats, clientAuthChain...)
	c.HandleFunc("POST /chats", c.CreateChat, clientAuthChain...)
	c.HandleFunc("PUT /chats/{chatID}/name", c.UpdateChatName, clientAuthChain...)

	// Участники
	c.HandleFunc("POST /chats/{chatID}/leave", c.LeaveChat, clientAuthChain...)
	c.HandleFunc("GET /chats/{chatID}/members", c.ChatMembers, clientAuthChain...)
	c.HandleFunc("DELETE /chats/{chatID}/members/{memberID}", c.DeleteMember, clientAuthChain...)

	// Приглашениями
	c.HandleFunc("GET /invitations", c.MyInvitations, clientAuthChain...)
	c.HandleFunc("GET /chats/{chatID}/invitations", c.ChatInvitations, clientAuthChain...)
	c.HandleFunc("POST /invitations", c.SendInvitation, clientAuthChain...)
	c.HandleFunc("POST /invitations/{invitationID}/accept", c.AcceptInvitation, clientAuthChain...)
	c.HandleFunc("POST /invitations/{invitationID}/cancel", c.CancelInvitation, clientAuthChain...)

}
