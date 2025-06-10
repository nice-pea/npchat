package register_handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
)

// Ping регистрирует обработчик для проверки работоспособности сервера.
// Данный обработчик не требует авторизации и может использоваться для health-check'а.
//
// Метод: GET /ping
func Ping(router http2.Router) {
	router.HandleFunc(
		"/ping",
		middleware.EmptyChain, // Пустая цепочка middleware, обработчик доступен без ограничений
		func(context http2.Context) (any, error) {
			// Возвращаем простую строку "pong" для подтверждения работоспособности сервера.
			return "pong", nil
		})
}
