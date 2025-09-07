package oauthAuthorize

import (
	"net/url"
	"testing"

	testifySuite "github.com/stretchr/testify/suite"

	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_OauthAuthorize() {
	usecase := &OauthAuthorizeUsecase{
		Providers: oauth.Providers{},
	}
	usecase.Providers.Add(suite.Adapters.Oauth)

	suite.Run("Provider обязательные поля, должен быть известен в сервисе", func() {
		out, err := usecase.OauthAuthorize(In{
			Provider:         "",
			CompleteCallback: "http://callback.ab",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)

		out, err = usecase.OauthAuthorize(In{
			Provider: "unknownProvider",
		})
		suite.ErrorIs(err, oauth.ErrUnknownOauthProvider)
		suite.Zero(out)
	})

	suite.Run("CompleteCallback должен быть корректным url", func() {
		out, err := usecase.OauthAuthorize(In{
			Provider:         suite.Adapters.Oauth.Name(),
			CompleteCallback: "",
		})
		suite.ErrorIs(err, ErrInvalidCompleteCallback)
		suite.Zero(out)

		out, err = usecase.OauthAuthorize(In{
			Provider:         suite.Adapters.Oauth.Name(),
			CompleteCallback: "adf[o",
		})
		suite.ErrorIs(err, ErrInvalidCompleteCallback)
		suite.Zero(out)
	})

	suite.Run("инициализация вернет валидный url", func() {
		// Инициализация регистрации
		out, err := usecase.OauthAuthorize(In{
			Provider:         suite.Adapters.Oauth.Name(),
			CompleteCallback: "http://callback.ab",
		})
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Вернется url
		parsedUrl, err := url.Parse(out.RedirectURL)
		suite.NoError(err)

		// Есть query-параметр state
		state := parsedUrl.Query().Get("state")
		suite.NotZero(state)
		// Есть query-параметр state
		code := parsedUrl.Query().Get("code")
		suite.NotZero(code)
		// Есть query-параметр redirect_uri
		ru := parsedUrl.Query().Get("redirect_uri")
		suite.NotZero(ru)
	})
}
