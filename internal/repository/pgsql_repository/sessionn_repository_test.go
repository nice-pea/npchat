package pgsqlRepository

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

func (suite *Suite) Test_SessionnRepository() {
	suite.Run("List", func() {
		suite.Run("из пустого репозитория вернется пустой список", func() {
			sessions, err := suite.RR.Sessions.List(sessionn.Filter{})
			suite.NoError(err)
			suite.Empty(sessions)
		})

		suite.Run("без фильтра из репозитория вернутся все сохраненные элементы", func() {
			sessions := suite.upsertRndSessions(10)
			fromRepo, err := suite.RR.Sessions.List(sessionn.Filter{})
			suite.NoError(err)
			suite.Len(fromRepo, len(sessions))
		})

		suite.Run("с фильтром по AccessToken вернутся, имеющие такой токен доступа", func() {
			// Создать много
			sessions := suite.upsertRndSessions(10)
			// Определить случайны искомый
			expected := common.RndElem(sessions)
			// Получить список
			fromRepo, err := suite.RR.Sessions.List(sessionn.Filter{
				AccessToken: expected.AccessToken.Token,
			})
			// Сравнить ожидания и результат
			suite.NoError(err)
			suite.Require().Len(fromRepo, 1)
			suite.Equal(expected, fromRepo[0])
		})
	})

	suite.Run("Upsert", func() {
		suite.Run("нельзя сохранять без ID", func() {
			err := suite.RR.Sessions.Upsert(sessionn.Session{
				ID:   uuid.Nil,
				Name: "someName",
			})
			suite.Error(err)
		})

		suite.Run("остальные поля, кроме ID могут быть пустыми", func() {
			err := suite.RR.Sessions.Upsert(sessionn.Session{
				ID: uuid.New(),
			})
			suite.NoError(err)
		})

		suite.Run("сохраненная сущность полностью соответствует сохраняемой", func() {
			// Создать и Сохранить
			session := suite.upsertSession(suite.rndSession())

			// Прочитать из репозитория
			sessions, err := suite.RR.Sessions.List(sessionn.Filter{})
			suite.NoError(err)
			suite.Require().Len(sessions, 1)
			suite.Equal(session, sessions[0])
		})

		suite.Run("перезапись с новыми значениями по ID", func() {
			id := uuid.New()
			// Несколько промежуточных состояний
			for range 33 {
				session := suite.rndSession()
				session.ID = id
				suite.upsertSession(session)
			}
			// Последнее сохраненное состояние
			expected := suite.rndSession()
			expected.ID = id
			suite.upsertSession(expected)

			// Прочитать из репозитория
			sessions, err := suite.RR.Sessions.List(sessionn.Filter{})
			suite.NoError(err)
			suite.Require().Len(sessions, 1)
			suite.Equal(expected, sessions[0])
		})
	})
}

func (suite *Suite) rndSession() sessionn.Session {
	suite.T().Helper()
	session, err := sessionn.NewSession(uuid.New(), gofakeit.UserAgent(), common.RndElem(sessionn.Statuses()))
	suite.Require().NoError(err)

	return session
}

func (suite *Suite) upsertRndSessions(count int) []sessionn.Session {
	suite.T().Helper()
	ss := make([]sessionn.Session, count)
	for i := range ss {
		ss[i] = suite.rndSession()
		suite.upsertSession(ss[i])
	}

	return ss
}

func (suite *Suite) upsertSession(session sessionn.Session) sessionn.Session {
	suite.T().Helper()
	err := suite.RR.Sessions.Upsert(session)
	suite.Require().NoError(err)

	return session
}
