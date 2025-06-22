package router

import (
	"net/http"

	"github.com/nice-pea/npchat/internal/controller/http2"
)

// Router обрабатывает HTTP-запросы
type Router struct {
	Services http2.Services
	http.ServeMux
}

func (c *Router) HandleFunc(pattern string, chain []http2.Middleware, handlerFunc http2.HandlerFunc) {
	handlerFuncRW := http2.WrapHandlerWithMiddlewares(handlerFunc, chain...)
	httpHandlerFunc := c.modulation(handlerFuncRW)
	c.ServeMux.HandleFunc(pattern, httpHandlerFunc)
	//slog.Info("Router: Зарегистрирован новый обработчик",
	//	slog.String("pattern", pattern),
	//)
}
