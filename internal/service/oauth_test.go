package service

import (
	"net/url"

	"github.com/google/uuid"
	"golang.org/x/exp/maps"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func (suite *servicesTestSuite) googleRegistrationInit() (out struct {
	out   GoogleRegistrationInitOut
	state string
}) {
	var err error
	out.out, err = suite.ss.oauth.GoogleRegistrationInit()
	suite.Require().NoError(err)
	suite.Require().NotZero(out)
	suite.Require().NotZero(out.out.RedirectURL)

	redirectUrl, err := url.Parse(out.out.RedirectURL)
	suite.Require().NoError(err)
	out.state = redirectUrl.Query().Get("state")
	suite.Require().NotEmpty(out.state)

	return out
}

func (suite *servicesTestSuite) Test_OAuth_GoogleRegistrationInit() {
	suite.Run("после инициализации, можно прочитать link из репозитория", func() {
		// Инициализация регистрации
		out, err := suite.ss.oauth.GoogleRegistrationInit()
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Вернется url с query-параметром state
		redirectUrl, err := url.Parse(out.RedirectURL)
		suite.NoError(err)
		state := redirectUrl.Query().Get("state")
		suite.Require().NotEmpty(state)

		// Прочитать из репозитория
		links, err := suite.rr.oauth.ListLinks(domain.OAuthListLinksFilter{})
		suite.NoError(err)
		suite.Require().Len(links, 1)
		suite.Equal(state, links[0].ID)
		suite.Empty(links[0].UserID)
		suite.Empty(links[0].ExternalID)
	})
}

func (suite *servicesTestSuite) Test_OAuth_GoogleRegistration() {
	suite.Run("UserCode обязательное поле", func() {
		input := GoogleRegistrationInput{
			UserCode:  "",
			InitState: uuid.NewString(),
		}
		out, err := suite.ss.oauth.GoogleRegistration(input)
		suite.ErrorIs(err, ErrInvalidUserCode)
		suite.Zero(out)
	})

	suite.Run("InitState обязательное поле", func() {
		input := GoogleRegistrationInput{
			UserCode:  uuid.NewString(),
			InitState: "",
		}
		out, err := suite.ss.oauth.GoogleRegistration(input)
		suite.ErrorIs(err, ErrInvalidInitState)
		suite.Zero(out)
	})

	suite.Run("неверный UserCode", func() {
		// Инициализация регистрации
		glOut := suite.googleRegistrationInit()

		input := GoogleRegistrationInput{
			UserCode:  uuid.NewString(),
			InitState: glOut.state,
		}
		out, err := suite.ss.oauth.GoogleRegistration(input)
		suite.ErrorIs(err, ErrWrongUserCode)
		suite.Zero(out)
	})

	suite.Run("неверный InitState", func() {
		input := GoogleRegistrationInput{
			UserCode:  maps.Keys(suite.mockOauthCodes)[0],
			InitState: uuid.NewString(),
		}
		out, err := suite.ss.oauth.GoogleRegistration(input)
		suite.ErrorIs(err, ErrWrongInitState)
		suite.Zero(out)
	})

	suite.Run("после регистрации пользователя можно прочитать", func() {
		// Инициализация регистрации
		glOut := suite.googleRegistrationInit()
		glCode := maps.Keys(suite.mockOauthCodes)[0]
		glToken := suite.mockOauthCodes[glCode]
		glUser := suite.mockGoogleUsers[glToken]

		// Завершить регистрацию
		input := GoogleRegistrationInput{
			UserCode:  glCode,
			InitState: glOut.state,
		}
		user, err := suite.ss.oauth.GoogleRegistration(input)
		suite.NoError(err)
		suite.Require().NotZero(user)

		// Проверить users репозиторий
		users, err := suite.rr.users.List(domain.UsersFilter{})
		suite.NoError(err)
		suite.Require().Len(users, 1)
		suite.Equal(user, users[0])
		// google
		suite.Equal(glUser.Name, user.Name)

		// Проверить oauth репозиторий
		links, err := suite.rr.oauth.ListLinks(domain.OAuthListLinksFilter{})
		suite.NoError(err)
		suite.Require().Len(links, 1)
		suite.Equal(glUser.ID, links[0].ExternalID)
		suite.Equal(user.ID, links[0].UserID)
		suite.Equal(glOut.state, links[0].ID)
	})
}
