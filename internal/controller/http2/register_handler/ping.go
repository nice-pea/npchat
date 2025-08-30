package register_handler

import (
	"github.com/gofiber/fiber/v2"
)

// Ping регистрирует обработчик для проверки работоспособности сервера.
// Данный обработчик не требует авторизации и может использоваться для health-check'а.
//
// Метод: GET /ping
func Ping(router *fiber.App) {
	router.Get(
		"/ping",
		func(context *fiber.Ctx) error {
			// Возвращаем простую строку "pong" для подтверждения работоспособности сервера.
			return context.SendString("pong")
		})
}
