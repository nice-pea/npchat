package repository_tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain/common"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

func TestRepository(t *testing.T, newRepository func() userr.Repository) {
	t.Run("List", func(t *testing.T) {
		t.Run("из пустого репозитория вернется пустой список", func(t *testing.T) {
			r := newRepository()
			users, err := r.List(userr.Filter{})
			assert.NoError(t, err)
			assert.Empty(t, users)
		})

		t.Run("без фильтра из репозитория вернутся все сохраненные элементы", func(t *testing.T) {
			r := newRepository()
			users := upsertRndUsers(t, r, 10)
			fromRepo, err := r.List(userr.Filter{})
			assert.NoError(t, err)
			assert.Len(t, fromRepo, len(users))
		})

		t.Run("с фильтром по ID вернется сохраненный элемент", func(t *testing.T) {
			r := newRepository()
			// Создать много
			users := upsertRndUsers(t, r, 10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := r.List(userr.Filter{
				ID: expected.ID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, fromRepo, 1)
			assert.Equal(t, expected, fromRepo[0])
		})

		t.Run("с фильтром по OAuthUserID вернутся, имеющие связь с пользователем oauth провайдера", func(t *testing.T) {
			r := newRepository()
			// Создать много
			users := upsertRndUsers(t, r, 10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := r.List(userr.Filter{
				OAuthUserID: expected.OpenAuthUsers[0].ID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, fromRepo, 1)
			assert.Equal(t, expected, fromRepo[0])
		})

		t.Run("с фильтром по OAuthProvider вернутся, имеющие связь с oauth провайдером", func(t *testing.T) {
			r := newRepository()
			// Создать много
			users := upsertRndUsers(t, r, 10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := r.List(userr.Filter{
				OAuthProvider: expected.OpenAuthUsers[0].Provider,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, fromRepo, 1)
			assert.Equal(t, expected, fromRepo[0])
		})

		t.Run("с фильтром по BasicAuthLogin вернутся, имеющие этот логин", func(t *testing.T) {
			r := newRepository()
			// Создать много
			users := upsertRndUsers(t, r, 10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := r.List(userr.Filter{
				BasicAuthLogin: expected.BasicAuth.Login,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, fromRepo, 1)
			assert.Equal(t, expected, fromRepo[0])
		})

		t.Run("с фильтром по BasicAuthPassword вернутся, имеющие этот пароль", func(t *testing.T) {
			r := newRepository()
			// Создать много
			users := upsertRndUsers(t, r, 10)
			// Определить случайны искомый
			expected := common.RndElem(users)
			// Получить список
			fromRepo, err := r.List(userr.Filter{
				BasicAuthPassword: expected.BasicAuth.Password,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, fromRepo, 1)
			assert.Equal(t, expected, fromRepo[0])
		})

		t.Run("можно искать по всем фильтрам сразу", func(t *testing.T) {
			r := newRepository()
			// Создать много чатов
			users := upsertRndUsers(t, r, 10)
			// Определить случайны искомый чат
			expected := common.RndElem(users)

			// Получить список
			fromRepo, err := r.List(userr.Filter{
				//ID:                expected.ID,
				OAuthUserID:       expected.OpenAuthUsers[0].ID,
				OAuthProvider:     expected.OpenAuthUsers[0].Provider,
				BasicAuthLogin:    expected.BasicAuth.Login,
				BasicAuthPassword: expected.BasicAuth.Password,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, fromRepo, 1)
			assert.Equal(t, expected, fromRepo[0])
		})
	})
	t.Run("Upsert", func(t *testing.T) {
		t.Run("нельзя сохранять без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(userr.User{
				ID:   "",
				Name: "someName",
			})
			assert.Error(t, err)
		})

		t.Run("остальные поля, кроме ID могут быть пустыми", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(userr.User{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
		})

		t.Run("сохраненная сущность полностью соответствует сохраняемой", func(t *testing.T) {
			r := newRepository()
			// Создать
			user := rndUser(t)
			addRndBasicAuth(t, &user)
			addRndOpenAuth(t, &user)

			// Сохранить
			err := r.Upsert(user)
			require.NoError(t, err)

			// Прочитать из репозитория
			users, err := r.List(userr.Filter{})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			assert.Equal(t, user, users[0])
		})

		t.Run("перезапись с новыми значениями по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			// Несколько промежуточных состояний
			for range 33 {
				user := rndUser(t)
				user.ID = id
				upsertUser(t, r, user)
			}
			// Последнее сохраненное состояние
			expected := rndUser(t)
			expected.ID = id
			upsertUser(t, r, expected)

			// Прочитать из репозитория
			users, err := r.List(userr.Filter{})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			assert.Equal(t, expected, users[0])
		})
	})
}

func rndUser(t *testing.T) userr.User {
	u, err := userr.NewUser(gofakeit.Name(), gofakeit.Noun())
	require.NoError(t, err)

	return u
}

func upsertRndUsers(t *testing.T, r userr.Repository, count int) []userr.User {
	users := make([]userr.User, count)
	for i := range users {
		users[i] = rndUser(t)
		addRndBasicAuth(t, &users[i])
		addRndOpenAuth(t, &users[i])
		upsertUser(t, r, users[i])
	}

	return users
}

func addRndBasicAuth(t *testing.T, user *userr.User) {
	ba, err := userr.NewBasicAuth(gofakeit.Noun(), common.RndPassword())
	require.NoError(t, err)
	err = user.AddBasicAuth(ba)
	require.NoError(t, err)
}

func addRndOpenAuth(t *testing.T, user *userr.User) {
	token, err := userr.NewOpenAuthToken(uuid.NewString(), "test", uuid.NewString(), time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	openAuthUser, err := userr.NewOpenAuthUser(uuid.NewString(), gofakeit.Company(), gofakeit.Email(), gofakeit.Noun(), gofakeit.URL(), token)
	require.NoError(t, err)
	err = user.AddOpenAuthUser(openAuthUser)
	require.NoError(t, err)
}

func upsertUser(t *testing.T, r userr.Repository, user userr.User) userr.User {
	err := r.Upsert(user)
	require.NoError(t, err)

	return user
}
