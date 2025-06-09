package register_handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// LoginByPassword регистрирует обработчик, позволяющий авторизоваться по логину и паролю.
// Доступен без предварительной аутентификации (публичная цепочка middleware).
//
// Метод: POST /auth/password/login
func LoginByPassword(router http2.Router) {
	// Тело запроса для авторизации по логину и паролю.
	type requestBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	router.HandleFunc(
		"POST /auth/password/login",
		middleware.ClientPubChain, // Цепочка middleware для публичных обработчиков (без авторизации)
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.BasicAuthLoginIn{
				Login:    rb.Login,
				Password: rb.Password,
			}

			return context.Services().Users().BasicAuthLogin(input)
		})
}

// RegistrationByPassword регистрирует обработчик, позволяющий регистрироваться по логину и паролю.
// Доступен без предварительной аутентификации (публичная цепочка middleware).
//
// Метод: POST /auth/password/registration
func RegistrationByPassword(router http2.Router) {
	// Тело запроса для авторизации по логину и паролю.
	type requestBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Nick     string `json:"nick"`
	}
	router.HandleFunc(
		"POST /auth/password/registration",
		middleware.ClientPubChain, // Цепочка middleware для публичных обработчиков (без авторизации)
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.BasicAuthRegistrationIn{
				Login:    rb.Login,
				Password: rb.Password,
				Name:     rb.Name,
				Nick:     rb.Nick,
			}

			return context.Services().Users().BasicAuthRegistration(input)
		})
}
