package http2

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

type Router interface {
	HandleFunc(pattern string, chain []Middleware, handler HandlerFunc)
}

type HandlerFunc func(Context) (any, error)

type Context interface {
	RequestID() string
	Session() domain.Session
	Request() *http.Request
	Writer() http.ResponseWriter
	Services() Services
}

type HandlerFuncRW func(RWContext) (any, error)

type RWContext interface {
	Context
	SetSession(domain.Session)
	SetRequestID(string)
	SetRequest(*http.Request)
}

type Middleware func(rw HandlerFuncRW) HandlerFuncRW

type Services interface {
	Chats() *service.Chats
	Invitations() *service.Invitations
	Members() *service.Members
	Sessions() *service.Sessions
	AuthnPassword() *service.AuthnPassword
	OAuth() *service.OAuth
}
