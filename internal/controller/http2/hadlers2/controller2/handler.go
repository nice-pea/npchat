package controller2

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/service"
)

type Handler2 interface {
	HandlerParams() HandlerParams
	HandlerFunc(Context) (any, error)
	//Middlewares() []Middleware
}

type HandlerParams struct {
	Path   string
	Method string
}

type Handler interface {
	Pattern() string
	HandlerFunc() HandlerFunc
	Middlewares() []Middleware
}

type HandlerFunc func(Context) (any, error)

type Context interface {
	//DecodeBody(any) error
	//PathValue(string) string
	//QueryValue(string) string

	//RequestID() string
	//Session() domain.Session
	Request() *http.Request
	Services() Services
}

type Middleware func(HandlerFunc) HandlerFunc

type Services interface {
	Chats() *service.Chats
	Invitations() *service.Invitations
	Members() *service.Members
	Sessions() *service.Sessions
	AuthnPassword() *service.AuthnPassword
}
