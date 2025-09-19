package pgsqlRepository

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

func (suite *Suite) Test_UserrRepository() {
	suite.Run("List", func() {
		suite.Run("из пустого репозитория вернется пустой список", func() {
			users, err := suite.RR.Users.List(userr.Filter{})
			suite.NoError(err)
			suite.Empty(users)
		})

		suite.Run("без фильтра из репозитория вернутся все сохраненные элементы", func() {
			users := suite.upsertRndUsers(10)
			fromRepo, err := suite.RR.Users.List(userr.Filter{})
			suite.NoError(err)
			suite.Len(fromRepo, len(users))
		})

		suite.Run("с фильтром по ID вернется сохраненный элемент", func() {
			// Создать много
			users := suite.upsertRndUsers(10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := suite.RR.Users.List(userr.Filter{
				ID: expected.ID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(fromRepo, 1)
			suite.Equal(expected, fromRepo[0])
		})

		suite.Run("с фильтром по OauthUserID вернутся, имеющие связь с пользователем oauth провайдера", func() {
			// Создать много
			users := suite.upsertRndUsers(10)
			// Определить случайны искомый
			expected := common.RndElem(users)

			// Получить список
			fromRepo, err := suite.RR.Users.List(userr.Filter{
				OauthUserID: expected.OpenAuthUsers[0].ID,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(fromRepo, 1)
			suite.Equal(expected, fromRepo[0])
		})

		suite.Run("с фильтром по OauthProvider вернутся, имеющие связь с oauth провайдером", func() {
			// Создать много
			users := suite.upsertRndUsers(10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := suite.RR.Users.List(userr.Filter{
				OauthProvider: expected.OpenAuthUsers[0].Provider,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(fromRepo, 1)
			suite.Equal(expected, fromRepo[0])
		})

		suite.Run("с фильтром по BasicAuthLogin вернутся, имеющие этот логин", func() {
			// Создать много
			users := suite.upsertRndUsers(10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := suite.RR.Users.List(userr.Filter{
				BasicAuthLogin: expected.BasicAuth.Login,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(fromRepo, 1)
			suite.Equal(expected, fromRepo[0])
		})

		suite.Run("с фильтром по BasicAuthPassword вернутся, имеющие этот пароль", func() {
			// Создать много
			users := suite.upsertRndUsers(10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := suite.RR.Users.List(userr.Filter{
				BasicAuthPassword: expected.BasicAuth.Password,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(fromRepo, 1)
			suite.Equal(expected, fromRepo[0])
		})

		suite.Run("можно искать по всем фильтрам сразу", func() {
			// Создать много
			users := suite.upsertRndUsers(10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := suite.RR.Users.List(userr.Filter{
				ID:                expected.ID,
				OauthUserID:       expected.OpenAuthUsers[0].ID,
				OauthProvider:     expected.OpenAuthUsers[0].Provider,
				BasicAuthLogin:    expected.BasicAuth.Login,
				BasicAuthPassword: expected.BasicAuth.Password,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(fromRepo, 1)
			suite.Equal(expected, fromRepo[0])
		})
	})

	suite.Run("Upsert", func() {
		suite.Run("нельзя сохранять без ID", func() {
			err := suite.RR.Users.Upsert(userr.User{
				ID:   uuid.Nil,
				Name: "someName",
			})
			suite.Error(err)
		})

		suite.Run("остальные поля, кроме ID могут быть пустыми", func() {
			err := suite.RR.Users.Upsert(userr.User{
				ID: uuid.New(),
			})
			suite.NoError(err)
		})

		suite.Run("сохраненная сущность полностью соответствует сохраняемой", func() {
			user := suite.rndUser()
			// Создать
			suite.addRndBasicAuth(&user)
			suite.addRndOpenAuth(&user)

			// Сохранить
			err := suite.RR.Users.Upsert(user)
			suite.Require().NoError(err)

			// Прочитать из репозитория
			users, err := suite.RR.Users.List(userr.Filter{})
			suite.NoError(err)
			suite.Require().Len(users, 1)
			suite.Equal(user, users[0])
		})

		suite.Run("перезапись с новыми значениями по ID", func() {
			id := uuid.New()
			// Несколько промежуточных состояний
			for range 33 {
				u := suite.rndUser()
				u.ID = id
				suite.upsertUser(u)
			}
			// Последнее сохраненное состояние
			expected := suite.rndUser()
			expected.ID = id
			suite.upsertUser(expected)

			// Прочитать из репозитория
			users, err := suite.RR.Users.List(userr.Filter{})
			suite.NoError(err)
			suite.Require().Len(users, 1)
			suite.Equal(expected, users[0])
		})
	})
}

// rndUser создает случайный экземпляр пользователя
func (suite *Suite) rndUser() userr.User {
	suite.T().Helper()
	u, err := userr.NewUser(gofakeit.Name(), gofakeit.Noun())
	suite.Require().NoError(err)
	return u
}

// upsertUser сохраняет пользователя в репозиторий
func (suite *Suite) upsertUser(user userr.User) userr.User {
	suite.T().Helper()
	err := suite.RR.Users.Upsert(user)
	suite.Require().NoError(err)
	return user
}

func (suite *Suite) upsertRndUsers(count int) []userr.User {
	suite.T().Helper()
	users := make([]userr.User, count)
	for i := range users {
		users[i] = suite.rndUser()
		suite.addRndBasicAuth(&users[i])
		suite.addRndOpenAuth(&users[i])
		suite.upsertUser(users[i])
	}

	return users
}

// addRndBasicAuth добавляет случайные базовые учетные данные пользователю
func (suite *Suite) addRndBasicAuth(user *userr.User) {
	suite.T().Helper()
	ba, err := userr.NewBasicAuth(gofakeit.Username(), common.RndPassword())
	suite.Require().NoError(err)
	err = user.AddBasicAuth(ba)
	suite.Require().NoError(err)
}

// addRndOpenAuth добавляет случайного oauth пользователя к пользователю
func (suite *Suite) addRndOpenAuth(user *userr.User) {
	suite.T().Helper()
	token, err := userr.NewOpenAuthToken(uuid.NewString(), "test", uuid.NewString(), time.Now().Add(1*time.Hour))
	suite.Require().NoError(err)
	openAuthUser, err := userr.NewOpenAuthUser(uuid.NewString(), gofakeit.Company(), gofakeit.Email(), gofakeit.Noun(), gofakeit.URL(), token)
	suite.Require().NoError(err)
	err = user.AddOpenAuthUser(openAuthUser)
	suite.Require().NoError(err)
}
