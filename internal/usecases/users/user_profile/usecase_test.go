package userProfile

import (
	"testing"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/userr"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) Test_UserProfile() {
	usecase := &UserProfileUsecase{
		Repo: suite.RR.Users,
	}
	suite.Run("есть валидация параметров", func() {
		// Некорректное значение subjectID
		_, err := usecase.UserProfile(In{UserID: uuid.New()})
		suite.ErrorIs(err, ErrInvalidSubjectID)
		// Некорректное значение userID
		_, err = usecase.UserProfile(In{SubjectID: uuid.New()})
		suite.ErrorIs(err, ErrInvalidUserID)
	})

	suite.Run("если пользователя с UserID не существует, вернется ошибка", func() {
		_, err := usecase.UserProfile(In{
			SubjectID: uuid.New(),
			UserID:    uuid.New(),
		})
		suite.ErrorIs(err, userr.ErrUserNotExists)
	})

	suite.Run("если userID равно subjectID, вернется полный профиль", func() {
		// Создаем пользователя
		user := suite.NewRndUserWithBasicAuth()
		user.OpenAuthUsers = []userr.OpenAuthUser{
			{ID: uuid.NewString()},
		}
		err := suite.RR.Users.Upsert(user)
		suite.Require().NoError(err)
		// Получаем профиль
		out, err := usecase.UserProfile(In{
			SubjectID: user.ID,
			UserID:    user.ID,
		})
		// Проверяем результат
		suite.NoError(err)
		suite.Equal(out.User, user)
	})

	suite.Run("если userID и subjectID не равны, вернется неполный профиль", func() {
		// Создаем пользователя
		user := suite.NewRndUserWithBasicAuth()
		user.OpenAuthUsers = []userr.OpenAuthUser{
			{ID: uuid.NewString()},
		}
		err := suite.RR.Users.Upsert(user)
		suite.Require().NoError(err)
		// Получаем профиль
		out, err := usecase.UserProfile(In{
			SubjectID: uuid.New(),
			UserID:    user.ID,
		})
		// Проверяем результат
		suite.NoError(err)
		suite.Equal(user.ID, out.User.ID)
		suite.Equal(user.Name, out.User.Name)
		suite.Equal(user.Nick, out.User.Nick)
		// Будут пустые поля
		suite.Empty(out.User.OpenAuthUsers)
		suite.Zero(out.User.BasicAuth)
	})
}
