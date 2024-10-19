package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/usecases"
)

type UserChats struct {
	UserChatsUc usecases.UserChatsUsecase
}

func (h *UserChats) Endpoint() string {
	return "/chats"
}

func (h *UserChats) Method() string {
	return http.MethodGet
}

func (h *UserChats) Fn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[UserChats] read body: %v", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		requestBody := _AuthRequestBody{}
		err = json.Unmarshal(b, &requestBody)
		if err != nil {
			log.Printf("[UserChats] Failed unmarshal request body: %v", err)
			http.Error(w, "Failed unmarshal request body", http.StatusBadRequest)
			return
		}
		out, err := h.UserChatsUc.UserChats(usecases.UserChatsIn{
			Login: requestBody.Login,
		})
		if err != nil {
			log.Printf("[UserChats] Failed handle healthcheck: %v", err)
			http.Error(w, "Failed handle healthcheck", http.StatusBadRequest)
			return
		}
		resp := _AuthResponse{
			AccessToken: out.AccessToken,
		}

		b, err = json.Marshal(resp)
		if err != nil {
			log.Printf("[UserChats] Failed marshal request body: %v", err)
			http.Error(w, "Failed marshal request body", http.StatusBadRequest)
			return
		}
		w.Write(b)
	}
}

type _UserChatsResponse struct {
	AccessToken string `json:"access_token"`
}

type _ChatApiModel struct {
	Id          int32             `json:"id"`
	Name        string            `json:"name"`
	LastMsg     *_MessageApiModel `json:"last_msg,omitempty"`
	LastReadMsg int32             `json:"last_read_msg,omitempty"`
	Unreadcount int               `json:"unread_count"`
	Permissions []string          `json:"permissions"`
}

type _MessageApiModel struct {
	Id         int32                    `json:"id"`
	ChatId     int32                    `json:"chat_id"`
	Date       int64                    `json:"date"`
	Text       string                   `json:"text,omitempty"`
	User       *_UserApiModel           `json:"user,omitempty"`
	Reply      *_ReplyedMessageApiModel `json:"reply,omitempty"`
	EditDate   int64                    `json:"edit_date,omitempty"`
	DeleteDate int64                    `json:"delete_date,omitempty"`
}

type _ReplyedMessageApiModel struct {
	Id         int32          `json:"id"`
	Date       int64          `json:"date"`
	Text       string         `json:"text,omitempty"`
	User       *_UserApiModel `json:"user,omitempty"`
	EditDate   int64          `json:"edit_date,omitempty"`
	DeleteDate int64          `json:"delete_date,omitempty"`
}
