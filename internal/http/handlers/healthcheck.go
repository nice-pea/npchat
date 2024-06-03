package handlers

import (
	"net/http"
)

type Healthcheck struct{}

func (h *Healthcheck) Endpoint() string {
	return "/healthcheck"
}

func (h *Healthcheck) Method() string {
	return http.MethodGet
}

func (h *Healthcheck) Fn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// b, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	logrus.Debug("[Healthcheck] read body: %v", err)
		// 	http.Error(w, "", http.StatusBadRequest)
		// 	return
		// }
		w.Write([]byte("ok"))
		w.WriteHeader(http.StatusOK)
	}
}

type _HealthcheckResponse struct {
	MinVersionSupport string `json: "min_version_support"`
	MinCodeSupport    int    `json: "min_code_support"`
	MaxVersionSupport string `json: "max_version_support"`
	MaxCodeSupport    int    `json: "max_code_support"`
}
