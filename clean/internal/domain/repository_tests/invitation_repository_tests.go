package repository_tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func InvitationsRepositoryTests(t *testing.T, newRepository func() domain.InvitationsRepository) {
	t.Helper()
	t.Run("List", func(t *testing.T) {
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			invitations, err := r.List(domain.InvitationsFilter{})
			assert.NoError(t, err)
			assert.Len(t, invitations, 0)
		})
		t.Run("без фильтра только один чат", func(t *testing.T) {
			r := newRepository()
			invitation := domain.Invitation{
				ID: uuid.NewString(),
			}
			err := r.Save(invitation)
			assert.NoError(t, err)
			invitations, err := r.List(domain.InvitationsFilter{})
			assert.NoError(t, err)
			assert.Len(t, invitations, 1)
		})
		t.Run("без фильтра отсутствует после удаления", func(t *testing.T) {
			r := newRepository()
			invitation := domain.Invitation{
				ID: uuid.NewString(),
			}
			err := r.Save(invitation)
			assert.NoError(t, err)
			invitations, err := r.List(domain.InvitationsFilter{})
			assert.NoError(t, err)
			assert.Len(t, invitations, 1)
			err = r.Delete(invitation.ID)
			assert.NoError(t, err)
			invitations, err = r.List(domain.InvitationsFilter{})
			assert.NoError(t, err)
			assert.Len(t, invitations, 0)
		})
		t.Run("с фильтром по id", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			assert.NoError(t, errors.Join(
				r.Save(domain.Invitation{ID: id}),
				r.Save(domain.Invitation{ID: uuid.NewString()}),
				r.Save(domain.Invitation{ID: uuid.NewString()}),
			))
			invitations, err := r.List(domain.InvitationsFilter{ID: id})
			assert.NoError(t, err)
			assert.Len(t, invitations, 1)
		})
		t.Run("с фильтром по chat id", func(t *testing.T) {
			r := newRepository()
			chat_id := uuid.NewString()
			assert.NoError(t, errors.Join(
				r.Save(domain.Invitation{ID: uuid.NewString(), ChatID: chat_id}),
				r.Save(domain.Invitation{ID: uuid.NewString(), ChatID: uuid.NewString()}),
				r.Save(domain.Invitation{ID: uuid.NewString(), ChatID: uuid.NewString()}),
			))
			invitations, err := r.List(domain.InvitationsFilter{ChatID: chat_id})
			assert.NoError(t, err)
			assert.Len(t, invitations, 1)
		})
		t.Run("с фильтром по user id", func(t *testing.T) {
			r := newRepository()
			const amountInvs = 2
			userID := uuid.NewString()
			invitations := make([]domain.Invitation, amountInvs)
			for i := range amountInvs {
				inv := domain.Invitation{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: userID, SubjectUserID: uuid.NewString()}
				invitations[i] = inv
				err := r.Save(inv)
				assert.NoError(t, err)
			}
			for range amountInvs {
				inv := domain.Invitation{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString(), SubjectUserID: uuid.NewString()}
				err := r.Save(inv)
				assert.NoError(t, err)
			}
			invitationsFromRepo, err := r.List(domain.InvitationsFilter{UserID: userID})
			assert.NoError(t, err)
			if assert.Len(t, invitationsFromRepo, amountInvs) {
				for i, inv := range invitationsFromRepo {
					assert.Equal(t, inv.ID, invitations[i].ID)
					assert.Equal(t, inv.ChatID, invitations[i].ChatID)
					assert.Equal(t, inv.UserID, invitations[i].UserID)
					assert.Equal(t, inv.SubjectUserID, invitations[i].SubjectUserID)
				}
			}
		})
		t.Run("с фильтром по subject user id", func(t *testing.T) {
			r := newRepository()
			const amountInvs = 2
			subjectUserID := uuid.NewString()
			invitations := make([]domain.Invitation, amountInvs)
			for i := range amountInvs {
				inv := domain.Invitation{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString(), SubjectUserID: subjectUserID}
				invitations[i] = inv
				err := r.Save(inv)
				assert.NoError(t, err)
			}
			for range amountInvs {
				inv := domain.Invitation{ID: uuid.NewString()}
				err := r.Save(inv)
				assert.NoError(t, err)
			}
			invitationsFromRepo, err := r.List(domain.InvitationsFilter{SubjectUserID: subjectUserID})
			assert.NoError(t, err)
			if assert.Len(t, invitationsFromRepo, amountInvs) {
				for i, inv := range invitationsFromRepo {
					assert.Equal(t, inv.ID, invitations[i].ID)
					assert.Equal(t, inv.ChatID, invitations[i].ChatID)
					assert.Equal(t, inv.UserID, invitations[i].UserID)
					assert.Equal(t, inv.SubjectUserID, invitations[i].SubjectUserID)
				}
			}
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("чат без id", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Invitation{
				ID: "",
			})
			assert.Error(t, err)
		})
		t.Run("чат без subject user id", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: "",
			})
			assert.NoError(t, err)
		})
		t.Run("чат без name", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Invitation{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Invitation{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
		})
		t.Run("дубль ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Invitation{
				ID: id,
			})
			assert.NoError(t, err)
			err = r.Save(domain.Invitation{
				ID: id,
			})
			assert.Error(t, err)
		})
		t.Run("дубль name", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Invitation{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
			err = r.Save(domain.Invitation{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
		})
		t.Run("дубль subject user id", func(t *testing.T) {
			r := newRepository()
			subjectUserId := uuid.NewString()
			err := r.Save(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: subjectUserId,
			})
			assert.NoError(t, err)
			err = r.Save(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: subjectUserId,
			})
			assert.NoError(t, err)
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("с пустым id", func(t *testing.T) {
			r := newRepository()
			err := r.Delete("")
			assert.Error(t, err)
		})
		t.Run("несуществующий id", func(t *testing.T) {
			r := newRepository()
			err := r.Delete(uuid.NewString())
			assert.NoError(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Invitation{
				ID: id,
			})
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
		t.Run("дважды удаленный", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Invitation{
				ID: id,
			})
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
	})
}
