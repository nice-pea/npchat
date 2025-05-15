package router

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
)

// Router обрабатывает HTTP-запросы
type Router struct {
	Services    http2.Services
	Middlewares []http2.Middleware
	http.ServeMux
}

//func InitRouter(sessions *service.Services, middlewares ...http2.Middleware) *Router {
//	c := &Router{
//		Services: sessions,
//		ServeMux: http.ServeMux{},
//
//	}
//
//	return c
//}

func (c *Router) HandleFunc(pattern string, handlerFunc http2.HandlerFunc) {
	handlerFunc = http2.Chain(handlerFunc, c.Middlewares...)

	c.ServeMux.HandleFunc(pattern, c.modulation(func(context http2.Context) (any, error) {
		return handlerFunc(context)
	}))
}
