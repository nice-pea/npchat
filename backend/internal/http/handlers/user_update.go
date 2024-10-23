package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/usecase"
)

type UserUpdate struct {
	UserUpdateUc usecase.UserUpdateUsecase
}

func (h *UserUpdate) Endpoint() string {
	return "/users"
}

func (h *UserUpdate) Method() string {
	return http.MethodPatch
}

func (h *UserUpdate) Fn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("[UserUpdate] read body: %v", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		requestBody := _UserUpdateRequest{}
		err = json.Unmarshal(b, &requestBody)
		if err != nil {
			log.Printf("[UserUpdate] Failed unmarshal request body: %v", err)
			http.Error(w, "Failed unmarshal request body", http.StatusBadRequest)
			return
		}
		_, err = h.UserUpdateUc.UserUpdate(usecase.UserUpdateIn{
			Username: requestBody.Username,
		})
		if err != nil {
			log.Printf("[UserUpdate] Failed handle UserByToken: %v", err)
			http.Error(w, "Failed handle UserByToken", http.StatusBadRequest)
			return
		}
		w.Write(b)
	}
}

type _UserUpdateRequest struct {
	Username string `json"username"`
}
