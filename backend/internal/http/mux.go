package http

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/service/l10n"
)

type mux struct {
	*http.ServeMux
	s ServerParams
}

type Request struct {
	*http.Request
	L10n   l10n.Service
	DB     *gorm.DB
	Locale string
}

func (m *mux) handle(pattern string, f func(Request) (any, error)) {
	m.Handle(pattern, wrap(m.s, f))
}

func wrap(s ServerParams, f func(Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := Request{
			Request: r,
			L10n:    s.L10n,
			Locale:  locale(r.Header.Get("Accept-Language"), l10n.LocaleDefault),
			DB:      s.DB,
		}
		var (
			data any
			b    []byte
			err  error
		)
		if err = r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			if data, err = f(req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				data = ResponseError{Error: err.Error()}
			}
			if s, ok := data.(string); ok {
				data = ResponseMsg{Message: s}
			}
		}
		// Marshal data
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
