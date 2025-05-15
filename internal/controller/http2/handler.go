package http2

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

type Router interface {
	HandleFunc(pattern string, handler HandlerFunc)
}

//	type Handler2 interface {
//		HandlerParams() HandlerParams
//		HandlerFunc(Context) (any, error)
//	}
//
//	type HandlerParams struct {
//		Path   string
//		Method string
//	}
//
//	type Handler interface {
//		Pattern() string
//		HandlerFunc() HandlerFunc
//		Middlewares() []Middleware
//	}
type HandlerFunc func(Context) (any, error)

type Context interface {
	//DecodeBody(any) error
	//PathValue(string) string
	//QueryValue(string) string

	RequestID() string
	Session() domain.Session
	Request() *http.Request
	Services() Services
}
type MutContext interface {
	//DecodeBody(any) error
	//PathValue(string) string
	//QueryValue(string) string

	//RequestID() string
	Context
	SetSession(domain.Session)
	SetRequestID(string)
	//SetServices() Services
}

type MiddlewareFunc func(MutContext) (any, error)
type Middleware func(MiddlewareFunc) MiddlewareFunc

type Services interface {
	Chats() *service.Chats
	Invitations() *service.Invitations
	Members() *service.Members
	Sessions() *service.Sessions
	AuthnPassword() *service.AuthnPassword
}
