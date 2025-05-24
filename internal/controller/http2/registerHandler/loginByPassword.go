package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// LoginByPassword регистрирует обработчик, позволяющий авторизоваться по логину и паролю.
// Доступен без предварительной аутентификации (публичная цепочка middleware).
//
// Метод: POST /login/password
func LoginByPassword(router http2.Router) {
	// requestBody описывает структуру тела запроса для авторизации по логину и паролю.
	type requestBody struct {
		Login    string `json:"login"`    // Логин пользователя
		Password string `json:"password"` // Пароль пользователя
	}
	router.HandleFunc(
		"POST /login/password",
		middleware.ClientPubChain, // Цепочка middleware для публичных обработчиков (без авторизации)
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			// Формируем входные данные для сервиса авторизации.
			input := service.AuthnPasswordLoginInput{
				Login:    rb.Login,
				Password: rb.Password,
			}

			// Вызываем сервис авторизации по логину и паролю и возвращаем результат.
			return context.Services().AuthnPassword().Login(input)
		})
}
