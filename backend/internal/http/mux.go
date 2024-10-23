package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/service/l10n"
)

type mux struct {
	*http.ServeMux
}

type Request struct {
	*http.Request
	L10n   l10n.Service
	Locale string
}

func (m *mux) handle(pattern string, f func(Request) (any, error)) {
	m.Handle(pattern, wrap(f))
}

func wrap(f func(Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{
			Request: r,
			L10n:    nil,
		}

		data, err := f(req)
		if err != nil {
			data = ResponseError{Error: err.Error()}
		}
		if s, ok := data.(string); ok {
			data = ResponseMsg{Message: s}
		}
		// Marshal data
		var b []byte
		if b, err = json.Marshal(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			b = []byte(err.Error())
		}
		// Try to send data
		if _, err = w.Write(b); err != nil {
			log.Println("[wrap] error response write:", err.Error())
			return
		}
	}
}

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}
