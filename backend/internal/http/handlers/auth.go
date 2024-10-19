package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/usecases"
)

type Auth struct {
	authUc usecases.AuthUsecase
}

func (h *Auth) Endpoint() string {
	return "/auth"
}

func (h *Auth) Method() string {
	return http.MethodPost
}

func (h *Auth) Fn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[Auth] read body: %v", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		requestBody := _AuthRequestBody{}
		err = json.Unmarshal(b, &requestBody)
		if err != nil {
			log.Printf("[Auth] Failed unmarshal request body: %v", err)
			http.Error(w, "Failed unmarshal request body", http.StatusBadRequest)
			return
		}
		out, err := h.authUc.Auth(usecases.AuthIn{
			Login: requestBody.Login,
		})
		if err != nil {
			log.Printf("[Auth] Failed handle healthcheck: %v", err)
			http.Error(w, "Failed handle healthcheck", http.StatusBadRequest)
			return
		}
		resp := _AuthResponse{
			AccessToken: out.AccessToken,
		}

		b, err = json.Marshal(resp)
		if err != nil {
			log.Printf("[Auth] Failed marshal request body: %v", err)
			http.Error(w, "Failed marshal request body", http.StatusBadRequest)
			return
		}
		w.Write(b)
	}
}

type _AuthRequestBody struct {
	Login string `json:"login"`
}

type _AuthResponse struct {
	AccessToken string `json:"access_token"`
}
