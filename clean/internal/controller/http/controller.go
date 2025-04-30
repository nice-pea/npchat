package http

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Context представляет контекст HTTP-запроса
type Context struct {
	requestID string
	subjectID string
	request   *http.Request
}

type HandlerFunc func(Context) (any, error)

// Controller обрабатывает HTTP-запросы
type Controller struct {
	chats       service.Chats
	invitations service.Invitations
	members     service.Members
	auth        interface {
		ByHeader(r *http.Request) (domain.User, error)
	}

	mux *http.ServeMux
}

func InitController(chats service.Chats, invitations service.Invitations, members service.Members) *Controller {
	c := &Controller{
		chats:       chats,
		invitations: invitations,
		members:     members,
		mux:         http.NewServeMux(),
	}
	registerHandlers(c)

	return c
}

// ServeHTTP обрабатывает HTTP-запросы
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.mux.ServeHTTP(w, r)
}
