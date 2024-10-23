package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type mux struct {
	*http.ServeMux
}

func (m *mux) handle(pattern string, f func(r *http.Request) (any, error)) {
	m.Handle(pattern, wrap(f))
}

func wrap(f func(r *http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := f(r)
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
