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
			users := make([]userr.User, 10)
			for i := range users {
				users[i] = upsertUser(t, r, rndUser(t))
			}
			userFromRepo, err := r.List(userr.Filter{})
			assert.NoError(t, err)
			assert.Len(t, userFromRepo, len(users))
		})

		t.Run("с фильтром по ID вернется сохраненный элемент", func(t *testing.T) {
			r := newRepository()
			// Создать много
			for range 10 {
				upsertUser(t, r, rndUser(t))
			}
			// Определить случайны искомый
			expectedUser := upsertUser(t, r, rndUser(t))
			// Получить список
			userFromRepo, err := r.List(userr.Filter{
				ID: expectedUser.ID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, userFromRepo, 1)
			assert.Equal(t, expectedUser, userFromRepo[0])
		})

		t.Run("с фильтром по InvitationID вернутся чаты, имеющие с приглашение с таким ID", func(t *testing.T) {
			r := newRepository()
			// Создать много
			users := make([]userr.User, 10)
			for i := range users {
				users[i] = rndUser(t)
				addRndBasicAuth(t, &users[i])
				addRndOpenAuth(t, &users[i])
				upsertUser(t, r, users[i])
			}
			// Определить случайны искомый
			expected := common.RndElem(users)

			// Получить список
			chatsFromRepo, err := r.List(userr.Filter{
				ID:                "",
				OAuthUserID:       "",
				OAuthProvider:     "",
				BasicAuthLogin:    "",
				BasicAuthPassword: "",
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, 1)
			assert.Equal(t, expected, chatsFromRepo[0])
		})

		t.Run("можно искать по всем фильтрам сразу", func(t *testing.T) {
			r := newRepository()
			// Создать много чатов
			chats := make([]userr.User, 10)
			for i := range chats {
				chats[i] = rndChat(t)
				addRndParticipant(t, &chats[i])
				addRndInv(t, &chats[i])
				upsertChat(t, r, chats[i])
			}
			// Определить случайны искомый чат
			expectedChat := common.RndElem(chats)

			// Получить список
			chatsFromRepo, err := r.List(userr.Filter{
				ID:                "",
				OAuthUserID:       "",
				OAuthProvider:     "",
				BasicAuthLogin:    "",
				BasicAuthPassword: "",
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, 1)
			assert.Equal(t, expectedChat, chatsFromRepo[0])
		})

		t.Run("можно вернуться несколько элементов", func(t *testing.T) {
			r := newRepository()
			// Участник, который есть во многих чатах
			rndp := rndParticipant(t)
			// Создать много чатов с искомым участником
			const expectedCount = 10
			for range expectedCount {
				user := rndChat(t)
				err := user.AddParticipant(rndp)
				require.NoError(t, err)
				upsertChat(t, r, user)
			}
			// Создать несколько других чатов
			for range 21 {
				upsertChat(t, r, rndChat(t))
			}
			// Получить список
			chatsFromRepo, err := r.List(userr.Filter{
				ParticipantID: rndp.UserID,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, chatsFromRepo, expectedCount)
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
			expectedUser := rndUser(t)
			expectedUser.ID = id
			upsertUser(t, r, expectedUser)

			// Прочитать из репозитория
			users, err := r.List(userr.Filter{})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			assert.Equal(t, expectedUser, users[0])
		})
	})
}

func rndUser(t *testing.T) userr.User {
	u, err := userr.NewUser(gofakeit.Name(), gofakeit.Noun())
	require.NoError(t, err)
	return u
}

func addRndBasicAuth(t *testing.T, user *userr.User) {
	ba, err := userr.NewBasicAuth(gofakeit.Noun(), "Passw0rd!")
	require.NoError(t, err)
	err = user.AddBasicAuth(ba)
	require.NoError(t, err)
}

func addRndOpenAuth(t *testing.T, user *userr.User) {
	token, err := userr.NewOpenAuthToken(uuid.NewString(), "test", uuid.NewString(), time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	openAuthUser, err := userr.NewOpenAuthUser(uuid.NewString(), gofakeit.Company(), "", gofakeit.Noun(), "", token)
	require.NoError(t, err)
	err = user.AddOpenAuthUser(openAuthUser)
	require.NoError(t, err)
}

func upsertUser(t *testing.T, r userr.Repository, user userr.User) userr.User {
	err := r.Upsert(user)
	require.NoError(t, err)

	return user
}
