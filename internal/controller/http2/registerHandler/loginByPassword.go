package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// LoginByPassword регистрирует обработчик, позволяющий Авторизация по логину/паролю
func LoginByPassword(router http2.Router) {
	type requestBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	router.HandleFunc(
		"POST /login/password",
		middleware.ClientPubChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.AuthnPasswordLoginInput{
				Login:    rb.Login,
				Password: rb.Password,
			}

			return context.Services().AuthnPassword().Login(input)
		})
}
