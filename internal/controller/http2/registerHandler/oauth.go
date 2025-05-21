package registerHandler

import (
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

func GoogleRegistration(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/google/registration",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			out, err := context.Services().OAuth().GoogleRegistrationInit()
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

func GoogleRegistrationCallback(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/google/registration/callback",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			input := service.GoogleRegistrationInput{
				UserCode:  http2.FormStr(context, "code"),
				InitState: http2.FormStr(context, "state"),
			}

			return context.Services().OAuth().GoogleRegistration(input)
		},
	)
}
