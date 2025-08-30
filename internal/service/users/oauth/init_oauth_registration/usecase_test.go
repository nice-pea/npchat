package initOAuthRegistration

import (
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	serviceSuite "github.com/nice-pea/npchat/internal/service/suite"
	"github.com/nice-pea/npchat/internal/service/users/oauth"
)

func (suite *testSuite) newRndUserWithSession(sessionStatus string) (out struct {
	User    userr.User
	Session sessionn.Session
}) {
	var err error
	out.User, err = userr.NewUser(gofakeit.Name(), "")
	suite.Require().NoError(err)
	err = suite.RR.Users.Upsert(out.User)
	suite.Require().NoError(err)

	out.Session, err = sessionn.NewSession(out.User.ID, gofakeit.ChromeUserAgent(), sessionStatus)
	suite.Require().NoError(err)
	err = suite.RR.Sessions.Upsert(out.Session)
	suite.Require().NoError(err)

	return
}

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_OAuth_InitRegistration() {
	usecase := &InitOAuthRegistrationUsecase{
		Providers: oauth.OAuthProviders{},
	}

	suite.Run("Provider обязательное поле", func() {
		// Инициализация регистрации
		out, err := usecase.InitOAuthRegistration(In{
			Provider: "",
		})
		suite.ErrorIs(err, ErrInvalidProvider)
		suite.Zero(out)
	})

	suite.Run("Provider должен быть известен в сервисе", func() {
		// Инициализация регистрации
		input := In{
			Provider: "unknownProvider",
		}
		out, err := usecase.InitOAuthRegistration(input)
		suite.ErrorIs(err, oauth.ErrUnknownOAuthProvider)
		suite.Zero(out)
	})

	suite.Run("инициализация вернет валидный url", func() {
		// Инициализация регистрации
		out, err := usecase.InitOAuthRegistration(In{
			Provider: suite.Adapters.Oauth.Name(),
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
	})
}
