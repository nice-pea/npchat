package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

// Ping регистрирует обработчик для проверки работоспособности сервера.
// Данный обработчик не требует авторизации и может использоваться для health-check'а.
//
// Метод: GET /ping
func Ping(router http2.Router) {
	router.HandleFunc(
		"/ping",
		middleware.BaseChain,
		func(context http2.Context) (any, error) {
			// Возвращаем простую строку "pong" для подтверждения работоспособности сервера.
			return "pong", nil
		})
}
