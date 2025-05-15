package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Авторизация по логину/паролю
func RegisterLoginByPasswordHandler(router http2.Router) {
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

			session, err := context.Services().AuthnPassword().Login(input)
			if err != nil {
				return nil, err
			}

			return session, nil
		})
}
