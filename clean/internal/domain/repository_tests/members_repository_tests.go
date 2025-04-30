package repository_tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func assertEqualMembers(t *testing.T, expected, actual domain.Member) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.ChatID, actual.ChatID)
}

func MembersRepositoryTests(t *testing.T, newRepository func() domain.MembersRepository) {
	t.Helper()
	t.Run("List", func(t *testing.T) {
		t.Run("из пустого репозитория вернется пустой список", func(t *testing.T) {
			r := newRepository()
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Empty(t, members)
		})
		t.Run("без фильтра из репозитория вернутся все сохраненные участники", func(t *testing.T) {
			r := newRepository()
			members := make([]domain.Member, 10)
			for i := range members {
				members[i] = domain.Member{
					ID:     uuid.NewString(),
					UserID: uuid.NewString(),
					ChatID: uuid.NewString(),
				}
				err := r.Save(members[i])
				require.NoError(t, err)
			}
			membersFromRepo, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Len(t, membersFromRepo, len(members))
		})
		t.Run("после удаления в репозитории будет отсутствовать этот участник", func(t *testing.T) {
			r := newRepository()
			member := domain.Member{ID: uuid.NewString()}
			err := r.Save(member)
			require.NoError(t, err)
			members, err := r.List(domain.MembersFilter{})
			require.NoError(t, err)
			require.Len(t, members, 1)
			err = r.Delete(member.ID)
			require.NoError(t, err)
			members, err = r.List(domain.MembersFilter{})
			require.NoError(t, err)
			require.Empty(t, members)
		})
		t.Run("с фильтром по ID вернется сохраненный участник", func(t *testing.T) {
			r := newRepository()
			for range 10 {
				err := r.Save(domain.Member{ID: uuid.NewString()})
				assert.NoError(t, err)
			}
			expectedMember := domain.Member{ID: uuid.NewString()}
			err := r.Save(expectedMember)
			require.NoError(t, err)
			members, err := r.List(domain.MembersFilter{ID: expectedMember.ID})
			require.NoError(t, err)
			require.Len(t, members, 1)
			assertEqualMembers(t, expectedMember, members[0])
		})
		t.Run("с фильтром по ChatID вернутся участники с равным ChatID", func(t *testing.T) {
			r := newRepository()
			chatID := uuid.NewString()
			require.NoError(t, errors.Join(
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: chatID}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: chatID}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString()}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString()}),
			))
			members, err := r.List(domain.MembersFilter{ChatID: chatID})
			assert.NoError(t, err)
			assert.Len(t, members, 2)
		})
		t.Run("с фильтром по UserID вернется несколько участников", func(t *testing.T) {
			r := newRepository()
			userID := uuid.NewString()
			for range 10 {
				member := domain.Member{ID: uuid.NewString(), UserID: uuid.NewString()}
				err := r.Save(member)
				assert.NoError(t, err)
			}
			expectedMembers := []domain.Member{
				{ID: uuid.NewString(), UserID: userID},
				{ID: uuid.NewString(), UserID: userID},
			}
			for _, member := range expectedMembers {
				err := r.Save(member)
				assert.NoError(t, err)
			}
			members, err := r.List(domain.MembersFilter{UserID: userID})
			assert.NoError(t, err)
			if assert.Len(t, members, len(expectedMembers)) {
				for i, member := range expectedMembers {
					assertEqualMembers(t, member, members[i])
				}
			}
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("нельзя сохранять участника без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Member{
				ID:     "",
				ChatID: uuid.NewString(),
				UserID: uuid.NewString(),
			})
			assert.Error(t, err)
		})
		t.Run("можно сохранять участника без ChatID", func(t *testing.T) {
			r := newRepository()
			member := domain.Member{
				ID:     uuid.NewString(),
				ChatID: "",
				UserID: uuid.NewString(),
			}
			err := r.Save(member)
			assert.NoError(t, err)
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			if assert.Len(t, members, 1) {
				assertEqualMembers(t, member, members[0])
			}
		})
		t.Run("можно сохранять участника без UserID", func(t *testing.T) {
			r := newRepository()
			member := domain.Member{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: "",
			}
			err := r.Save(member)
			assert.NoError(t, err)
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			require.Len(t, members, 1)
			assertEqualMembers(t, member, members[0])
		})
		t.Run("сохраненный участник полностью соответствует сохраняемому", func(t *testing.T) {
			r := newRepository()
			member := domain.Member{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: uuid.NewString(),
			}
			err := r.Save(member)
			assert.NoError(t, err)
			members, err := r.List(domain.MembersFilter{})
			require.NoError(t, err)
			require.Len(t, members, 1)
			assertEqualMembers(t, member, members[0])
		})
		t.Run("перезапись с новыми значениями по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			const count = 33
			for range count {
				err := r.Save(domain.Member{
					ID:     id,
					UserID: uuid.NewString(),
					ChatID: uuid.NewString(),
				})
				assert.NoError(t, err)
			}
			expectedMember := domain.Member{
				ID:     id,
				UserID: uuid.NewString(),
				ChatID: uuid.NewString(),
			}
			err := r.Save(expectedMember)
			assert.NoError(t, err)
			members, err := r.List(domain.MembersFilter{ID: id})
			require.NoError(t, err)
			require.Len(t, members, 1)
			assertEqualMembers(t, expectedMember, members[0])
		})
		t.Run("ChatID может дублироваться", func(t *testing.T) {
			r := newRepository()
			chatID := uuid.NewString()
			const count = 33
			for range count {
				err := r.Save(domain.Member{
					ID:     uuid.NewString(),
					ChatID: chatID,
				})
				assert.NoError(t, err)
			}
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Len(t, members, count)
		})
		t.Run("UserID может дублироваться", func(t *testing.T) {
			r := newRepository()
			userID := uuid.NewString()
			const count = 33
			for range count {
				err := r.Save(domain.Member{
					ID:     uuid.NewString(),
					UserID: userID,
				})
				assert.NoError(t, err)
			}
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Len(t, members, count)
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("параметр ID обязательный", func(t *testing.T) {
			r := newRepository()
			err := r.Delete("")
			assert.Error(t, err)
		})
		t.Run("несуществующий ID не вернет ошибку", func(t *testing.T) {
			r := newRepository()
			err := r.Delete(uuid.NewString())
			assert.NoError(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Member{ID: id})
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
		t.Run("можно повторно удалять по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Member{ID: id})
			assert.NoError(t, err)
			for range 343 {
				err = r.Delete(id)
				assert.NoError(t, err)
			}
			chats, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Empty(t, chats)
		})
	})
}
