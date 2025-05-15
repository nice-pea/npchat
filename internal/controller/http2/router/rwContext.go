package router

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Context представляет контекст HTTP-запроса
type rwContext struct {
	requestID string
	request   *http.Request
	session   domain.Session
	services  http2.Services
}

func (r *rwContext) RequestID() string {
	return r.requestID
}

func (r *rwContext) Session() domain.Session {
	return r.session
}

func (r *rwContext) Request() *http.Request {
	return r.request
}

func (r *rwContext) Services() http2.Services {
	return r.services
}

func (r *rwContext) SetSession(session domain.Session) {
	r.session = session
}

func (r *rwContext) SetRequestID(id string) {
	r.requestID = id
}
