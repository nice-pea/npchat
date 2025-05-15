package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
)

// Ping регистрирует обработчик для проверки работоспособности сервера
func Ping(router http2.Router) {
	router.HandleFunc(
		"/ping",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			return "pong", nil
		})
}
