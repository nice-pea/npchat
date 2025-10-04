package redisRegistry_test

import (
	"time"

	redisRegistry "github.com/nice-pea/npchat/internal/adapter/jwt/registry/redis"

	"github.com/google/uuid"
)

// Test_Registry тестирование методов Registry для работы с Redis
func (suite *testSuite) Test_Registry() {
	suite.Run("RegisterIssueTime", func() {
		suite.Run("если sessionID пустой то вернется ошибка, в редис ничего не запишется", func() {
			err := suite.RedisCli.RegisterIssueTime(uuid.UUID{}, time.Now())
			suite.ErrorIs(err, redisRegistry.ErrEmptySessionID)
			suite.requireIsRedisEmpty()
		})

		suite.Run("если issueTime пустой то вернется ошибка, в редис ничего не запишется", func() {
			err := suite.RedisCli.RegisterIssueTime(uuid.New(), time.Time{})
			suite.ErrorIs(err, redisRegistry.ErrEmptyIssueTime)
			suite.requireIsRedisEmpty()
		})

		suite.Run("новое значение сменит старое, в редис будет только одна запись", func() {
			sessionId := uuid.New()
			issueTime := time.Now()
			err := suite.RedisCli.RegisterIssueTime(sessionId, issueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis1 := suite.getIssueTime(sessionId)
			suite.Require().Equal(issueTime, issueTimeFromRedis1)

			newIssueTime := time.Now().Add(time.Hour)
			err = suite.RedisCli.RegisterIssueTime(sessionId, newIssueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis2 := suite.getIssueTime(sessionId)
			suite.Require().Equal(newIssueTime, issueTimeFromRedis2)

			suite.Require().NotEqual(issueTimeFromRedis2, issueTimeFromRedis1)

			keys := suite.redisKeys()
			suite.Require().Len(keys, 1)
			suite.Equal(sessionId.String(), keys[0])
		})

		suite.Run("можно записать более одной записи", func() {
			for range 10 {
				err := suite.RedisCli.RegisterIssueTime(uuid.New(), time.Now())
				suite.Require().NoError(err)
			}

			keys := suite.redisKeys()
			suite.Require().Len(keys, 10)
		})

		suite.Run("созданная запись живет столкьоже сколько передали в поле Ttl", func() {
			id := uuid.New()
			// ставим Ttl в 1 миллисекунду
			suite.RedisCli.Ttl = time.Millisecond
			err := suite.RedisCli.RegisterIssueTime(id, time.Now())
			suite.Require().NoError(err)

			time.Sleep(10 * time.Millisecond)
			issueTime := suite.getIssueTime(id)
			suite.Require().Zero(issueTime)

			suite.requireIsRedisEmpty()
		})

		suite.Run("все записи которые Ttl прошел будут удалены из кэша", func() {
			suite.RedisCli.Ttl = time.Millisecond
			for range 10 {
				err := suite.RedisCli.RegisterIssueTime(uuid.New(), time.Now())
				suite.Require().NoError(err)
			}

			suite.RedisCli.Ttl = time.Minute
			id := uuid.New()
			err := suite.RedisCli.RegisterIssueTime(id, time.Now())
			suite.Require().NoError(err)

			time.Sleep(10 * time.Millisecond)

			keys := suite.redisKeys()
			suite.Require().Len(keys, 1)
			suite.Equal(id.String(), keys[0])
		})

	})
	suite.Run("IssueTime", func() {
		suite.Run("если sessionID пустой то вернется ошибка", func() {
			issueTime, err := suite.RedisCli.IssueTime(uuid.UUID{})
			suite.Require().ErrorIs(err, redisRegistry.ErrEmptySessionID)
			suite.Zero(issueTime)
		})

		suite.Run("если с таким sessionID нету значения в кэше то вернется ZeroValue, nil", func() {
			issueTime, err := suite.RedisCli.IssueTime(uuid.New())
			suite.Require().NoError(err)
			suite.Zero(issueTime)
		})

		suite.Run("если существует в кэше такой sessionId то вернется его значение", func() {
			sessionId := uuid.New()
			issueTime := time.Now()
			suite.setIssueTime(sessionId, issueTime, time.Minute)

			issueTimeRepo, err := suite.RedisCli.IssueTime(sessionId)
			suite.Require().NoError(err)
			suite.Equal(issueTime, issueTimeRepo)
		})

	})
}
