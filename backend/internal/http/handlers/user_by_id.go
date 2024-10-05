package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/saime-0/cute-chat-backend/internal/model"
	"github.com/saime-0/cute-chat-backend/internal/usecases"
	"github.com/sirupsen/logrus"
)

type UserByID struct {
	UserByIDUc usecases.UserByIDUsecase
}

func (h *UserByID) Endpoint() string {
	return "/users/:id"
}

func (h *UserByID) Method() string {
	return http.MethodGet
}

func (h *UserByID) Fn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			logrus.Debug("[UserByID] Path value `id` not received")
			http.Error(w, "Path value `id` not received", http.StatusBadRequest)
			return
		}
		out, err := h.UserByIDUc.UserByID(usecases.UserByIDIn{
			ID: model.ID(id),
		})
		if err != nil {
			logrus.Debugf("[UserByID] Failed handle UserByID: %v", err)
			http.Error(w, "Failed handle UserByID", http.StatusBadRequest)
			return
		}
		if !out.Found {
			logrus.Debug("[UserByID] User not found by id")
			http.Error(w, "User not found by id", http.StatusBadRequest)
			return
		}
		resp := _UserByIDResponse{
			User: userToApiModel(out.User),
		}
		b, err := json.Marshal(resp)
		if err != nil {
			logrus.Debugf("[UserByID] Failed marshal request body: %v", err)
			http.Error(w, "Failed marshal request body", http.StatusBadRequest)
			return
		}
		w.Write(b)
	}
}

type _UserByIDResponse struct {
	User _UserApiModel `json:"user"`
}
