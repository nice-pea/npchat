package redisRegistry_test

import (
	"time"

	redisRegistry "github.com/nice-pea/npchat/internal/adapter/jwt/registry/redis"

	"github.com/google/uuid"
)

// Test_Init тестовый сценарий для проверки функции инициализации Registry
func (suite *testSuite) Test_Init() {
	suite.Run("создание Registry с валидной конфигурацией", func() {
		cfg := redisRegistry.Config{
			DSN: suite.DSN,
			Ttl: 2 * time.Minute,
		}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().NoError(err)
		suite.NotNil(cli)
	})

	suite.Run("создание Registry с пустой конфигурацией", func() {
		cfg := redisRegistry.Config{}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().Error(err)
		suite.Zero(cli)
	})

	suite.Run("создание Registry с невалидной конфигурацией", func() {
		cfg := redisRegistry.Config{
			DSN: "241421.46334.14241.61253:253532325",
		}
		cli, err := redisRegistry.Init(cfg)

		suite.Require().Error(err)
		suite.Zero(cli)
	})
}

// Test_Registry тестирование методов Registry для работы с Redis
func (suite *testSuite) Test_Registry() {
	suite.Run("RegisterIssueTime", func() {
		suite.Run("если sessionID пустой то вернется ошибка, в редис ничего не запишется", func() {
			err := suite.Registry.RegisterIssueTime(uuid.UUID{}, time.Now())
			suite.ErrorIs(err, redisRegistry.ErrEmptySessionID)
			suite.requireIsRedisEmpty()
		})

		suite.Run("если issueTime пустой то вернется ошибка, в редис ничего не запишется", func() {
			err := suite.Registry.RegisterIssueTime(uuid.New(), time.Time{})
			suite.ErrorIs(err, redisRegistry.ErrEmptyIssueTime)
			suite.requireIsRedisEmpty()
		})

		suite.Run("новое значение сменит старое, в редис будет только одна запись", func() {
			sessionId := uuid.New()
			issueTime := time.Now()
			err := suite.Registry.RegisterIssueTime(sessionId, issueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis1 := suite.getIssueTime(sessionId)
			suite.Require().Equal(issueTime.Unix(), issueTimeFromRedis1.Unix())

			newIssueTime := time.Now().Add(time.Hour)
			err = suite.Registry.RegisterIssueTime(sessionId, newIssueTime)
			suite.Require().NoError(err)

			issueTimeFromRedis2 := suite.getIssueTime(sessionId)
			suite.Require().Equal(newIssueTime.Unix(), issueTimeFromRedis2.Unix())

			suite.Require().NotEqual(issueTimeFromRedis2.Unix(), issueTimeFromRedis1.Unix())

			keys := suite.redisKeys()
			suite.Require().Len(keys, 1)
			suite.Equal(sessionId.String(), keys[0])
		})

		suite.Run("можно записать более одной записи", func() {
			for range 10 {
				err := suite.Registry.RegisterIssueTime(uuid.New(), time.Now())
				suite.Require().NoError(err)
			}

			keys := suite.redisKeys()
			suite.Require().Len(keys, 10)
		})

		suite.Run("созданная запись живет столько же сколько передали в поле Ttl", func() {
			id := uuid.New()
			// ставим Ttl в 1 миллисекунду
			suite.Registry.Ttl = time.Millisecond
			err := suite.Registry.RegisterIssueTime(id, time.Now())
			suite.Require().NoError(err)

			time.Sleep(10 * time.Millisecond)
			issueTime := suite.getIssueTime(id)
			suite.Require().Zero(issueTime)

			suite.requireIsRedisEmpty()
		})

		suite.Run("все записи которые Ttl прошел будут удалены из кэша", func() {
			suite.Registry.Ttl = time.Millisecond
			for range 10 {
				err := suite.Registry.RegisterIssueTime(uuid.New(), time.Now())
				suite.Require().NoError(err)
			}
			time.Sleep(10 * time.Millisecond)
			keys := suite.redisKeys()
			suite.Require().Len(keys, 0)
		})

	})
	suite.Run("IssueTime", func() {
		suite.Run("если sessionID пустой то вернется ошибка", func() {
			issueTime, err := suite.Registry.IssueTime(uuid.UUID{})
			suite.Require().ErrorIs(err, redisRegistry.ErrEmptySessionID)
			suite.Zero(issueTime)
		})

		suite.Run("если с таким sessionID нету значения в кэше то вернется ZeroValue, nil", func() {
			issueTime, err := suite.Registry.IssueTime(uuid.New())
			suite.Require().NoError(err)
			suite.Zero(issueTime)
		})

		suite.Run("если существует в кэше такой sessionId то вернется его значение", func() {
			sessionId := uuid.New()
			issueTime := time.Now()
			suite.setIssueTime(sessionId, issueTime, time.Minute)

			issueTimeRepo, err := suite.Registry.IssueTime(sessionId)
			suite.Require().NoError(err)
			suite.Equal(issueTime.Unix(), issueTimeRepo.Unix())
		})

	})
}
