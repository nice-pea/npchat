package router

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

// Router обрабатывает HTTP-запросы
type Router struct {
	//chats         *service.Chats
	//invitations   *service.Invitations
	//members       *service.Members
	sessions *service.Sessions
	//authnPassword *service.AuthnPassword

	http.ServeMux
}

func InitRouter(sessions *service.Sessions) *Router {
	c := &Router{
		sessions: sessions,
		//chats:         chats,
		//invitations:   invitations,
		//members:       members,
		//authnPassword: authnPassword,
		ServeMux: http.ServeMux{},
	}
	//c.registerHandlers()

	return c
}

func (c *Router) HandleFunc(pattern string, handlerFunc HandlerFunc, middlewares ...middleware) {
	c.ServeMux.HandleFunc(pattern, c.modulation(chain(handlerFunc, middlewares...)))
}

type middleware func(HandlerFunc) HandlerFunc

func chain(h HandlerFunc, middlewares ...middleware) HandlerFunc {
	for _, mw := range slices.Backward(middlewares) {
		h = mw(h)
	}
	return h
}
