package redisCache_test

import (
	"time"

	redisCache "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"

	"github.com/google/uuid"
)

func (suite *testSuite) Test_JWTIssuanceRegistry() {
	suite.Run("RegisterIssueTime", func() {
		// RegisterIssueTime(sessionID uuid.UUID, issueTime time.Time) error {
		suite.Run("если sessionID пустой то вернется ошибка, в редис ничего не запишется", func() {
			err := suite.RedisCli.RegisterIssueTime(uuid.UUID{}, time.Now())
			suite.Assert().ErrorIs(err, redisCache.ErrEmptySessionID)
			suite.redisEmpty()
		})
		suite.Run("если issueTime пустой то вернется ошибка, в редис ничего не запишется", func() {
			err := suite.RedisCli.RegisterIssueTime(uuid.New(), time.Time{})
			suite.Assert().ErrorIs(err, redisCache.ErrEmptyIssueTime)
			suite.redisEmpty()
		})

		suite.Run("новое значение сменит старое, в редис будет только одна запись", func() {
			sessionId := uuid.New()
			issueTime := time.Now()
			err := suite.RedisCli.RegisterIssueTime(sessionId, issueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis1 := suite.getIssueTime(sessionId)
			suite.Require().True(issueTime.Equal(issueTimeFromRedis1))

			newIssueTime := time.Now().Add(time.Hour)
			err = suite.RedisCli.RegisterIssueTime(sessionId, newIssueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis2 := suite.getIssueTime(sessionId)
			suite.Require().True(newIssueTime.Equal(issueTimeFromRedis2))

			suite.Require().NotEqual(issueTimeFromRedis2, issueTimeFromRedis1)

			keys := suite.redisKeys()
			suite.Require().Len(keys, 1)
			suite.Assert().Equal(sessionId.String(), keys[0])

		})
		suite.Run("можно записать более одного записей", func() {
			for range 10 {
				err := suite.RedisCli.RegisterIssueTime(uuid.New(), time.Now())
				suite.Require().NoError(err)
			}

			keys := suite.redisKeys()
			suite.Require().Len(keys, 10)

		})

	})
	suite.Run("GetIssueTime", func() {
		// GetIssueTime(sessionID uuid.UUID) (*time.Time, error)
		suite.Run("из пустого репозитория вернется NULL", func() {

		})

	})
}
