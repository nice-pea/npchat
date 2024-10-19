package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/saime-0/nice-pea-chat/internal/model"
	"github.com/saime-0/nice-pea-chat/internal/usecases"
)

type UserByToken struct {
	UserByTokenUc usecases.UserByTokenUsecase
}

func (h *UserByToken) Endpoint() string {
	return "/users"
}

func (h *UserByToken) Method() string {
	return http.MethodGet
}

const _AUTH_HEADER = "Authorization"

func (h *UserByToken) Fn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(_AUTH_HEADER)
		token := strings.TrimSuffix(header, "Bearer ")
		if len(token) == 0 {
			log.Println("[UserByToken] AuthHeader is empty")
			http.Error(w, "AuthHeader is empty", http.StatusBadRequest)
			return
		}
		out, err := h.UserByTokenUc.UserByToken(usecases.UserByTokenIn{
			Token: token,
		})
		if err != nil {
			log.Printf("[UserByToken] Failed handle UserByToken: %v", err)
			http.Error(w, "Failed handle UserByToken", http.StatusBadRequest)
			return
		}
		if !out.Found {
			log.Println("[UserByToken] User not found by token")
			http.Error(w, "User not found by token", http.StatusBadRequest)
			return
		}
		resp := _UserByTokenResponse{
			User:  userToApiModel(out.User),
			Creds: credsToApiModel(out.Creds),
		}
		b, err := json.Marshal(resp)
		if err != nil {
			log.Printf("[UserByToken] Failed marshal request body: %v", err)
			http.Error(w, "Failed marshal request body", http.StatusBadRequest)
			return
		}
		w.Write(b)
	}
}

func userToApiModel(u model.User) _UserApiModel {
	return _UserApiModel{
		ID:       string(u.ID),
		Username: u.Username,
	}
}

func credsToApiModel(c model.Credentials) _CredsApiModel {
	return _CredsApiModel{
		Login: c.Login,
	}
}

type _UserByTokenResponse struct {
	User  _UserApiModel  `json:"user"`
	Creds _CredsApiModel `json:"credentials"`
}

type _UserApiModel struct {
	ID       string `json"id"`
	Username string `json"username"`
}

type _CredsApiModel struct {
	Login string `json:"login"`
}
