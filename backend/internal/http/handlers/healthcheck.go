package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/usecase"
)

type Healthcheck struct {
	HealthcheckUc usecase.HealthcheckUsecase
}

func (h *Healthcheck) Endpoint() string {
	return "/healthcheck"
}

func (h *Healthcheck) Method() string {
	return http.MethodGet
}

func (h *Healthcheck) Fn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := h.HealthcheckUc.Healthcheck()
		if err != nil {
			log.Printf("[Healthcheck] Failed handle healthcheck: %v", err)
			http.Error(w, "Failed handle healthcheck", http.StatusBadRequest)
			return
		}
		resp := _HealthcheckResponse(out)
		b, err := json.Marshal(resp)
		if err != nil {
			log.Printf("[Healthcheck] Failed marshal request body: %v", err)
			http.Error(w, "Failed marshal request body", http.StatusBadRequest)
			return
		}
		w.Write(b)
	}
}

type _HealthcheckResponse struct {
	MinVersionSupport string `json:"min_version_support"`
	MinCodeSupport    int    `json:"min_code_support"`
	MaxVersionSupport string `json:"max_version_support"`
	MaxCodeSupport    int    `json:"max_code_support"`
}
