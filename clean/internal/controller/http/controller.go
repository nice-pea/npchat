package http

import (
	"net/http"
	"slices"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Context представляет контекст HTTP-запроса
type Context struct {
	requestID string
	request   *http.Request
	session   domain.Session
}

type HandlerFunc func(Context) (any, error)

// Controller обрабатывает HTTP-запросы
type Controller struct {
	chats            *service.Chats
	invitations      *service.Invitations
	members          *service.Members
	sessions         *service.Sessions
	loginCredentials *service.LoginCredentials

	http.ServeMux
}

func InitController(chats *service.Chats, invitations *service.Invitations, members *service.Members, sessions *service.Sessions, loginCredentials *service.LoginCredentials) *Controller {
	c := &Controller{
		chats:            chats,
		invitations:      invitations,
		members:          members,
		sessions:         sessions,
		loginCredentials: loginCredentials,
		ServeMux:         http.ServeMux{},
	}
	c.registerHandlers()

	return c
}

func (c *Controller) HandleFunc(pattern string, handlerFunc HandlerFunc, middlewares ...middleware) {
	c.ServeMux.HandleFunc(pattern, c.modulation(chain(handlerFunc, middlewares...)))
}

type middleware func(HandlerFunc) HandlerFunc

func chain(h HandlerFunc, middlewares ...middleware) HandlerFunc {
	for _, mw := range slices.Backward(middlewares) {
		h = mw(h)
	}
	return h
}
