package registerHandler

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

func OAuthRegistration(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/registration",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			input := service.OAuthInitRegistrationInput{
				Provider: http2.PathStr(context, "provider"),
			}
			out, err := context.Services().OAuth().InitRegistration(input)
			if err != nil {
				return nil, err
			}

			return http2.Redirect{
				URL:  out.RedirectURL,
				Code: http.StatusTemporaryRedirect,
			}, nil
		},
	)
}

func OAuthRegistrationCallback(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/registration/callback",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			input := service.OAuthCompeteRegistrationInput{
				UserCode:  http2.FormStr(context, "code"),
				InitState: http2.FormStr(context, "state"),
				Provider:  http2.PathStr(context, "provider"),
			}

			return context.Services().OAuth().CompeteRegistration(input)
		},
	)
}
