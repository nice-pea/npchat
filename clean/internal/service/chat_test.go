package service

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	domain_mock "github.com/saime-0/nice-pea-chat/internal/domain/mocks"
)

func TestChats_Create(t *testing.T) {
	t.Run("успешное созданиие в моковом репозитории", func(t *testing.T) {
		// Настройка репозитория
		repo := new(domain_mock.ChatsRepositoryMock)
		repo.On("Save", mock.Anything).Return(nil)
		service := &Chats{
			ChatsRepo: repo,
			History:   HistoryDummy{},
		}
		// Тест функции
		assert.NoError(t, service.Create(domain.Chat{
			ID:   uint(rand.Int()),
			Name: "abcd",
		}))
	})
}

func TestChats_Delete(t *testing.T) {
	type fields struct {
		ChatsRepo domain.ChatsRepository
		History   History
	}
	type args struct {
		memberID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &Chats{
				ChatsRepo: tt.fields.ChatsRepo,
				History:   tt.fields.History,
			}
			tt.wantErr(t, ch.Delete(tt.args.memberID), fmt.Sprintf("Delete(%v)", tt.args.memberID))
		})
	}
}

func TestChats_List(t *testing.T) {
	type fields struct {
		ChatsRepo domain.ChatsRepository
		History   History
	}
	type args struct {
		memberID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &Chats{
				ChatsRepo: tt.fields.ChatsRepo,
				History:   tt.fields.History,
			}
			tt.wantErr(t, ch.List(tt.args.memberID), fmt.Sprintf("List(%v)", tt.args.memberID))
		})
	}
}

func TestChats_Members(t *testing.T) {
	type fields struct {
		ChatsRepo domain.ChatsRepository
		History   History
	}
	type args struct {
		chatID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Member
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &Chats{
				ChatsRepo: tt.fields.ChatsRepo,
				History:   tt.fields.History,
			}
			got, err := ch.Members(tt.args.chatID)
			if !tt.wantErr(t, err, fmt.Sprintf("Members(%v)", tt.args.chatID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Members(%v)", tt.args.chatID)
		})
	}
}
