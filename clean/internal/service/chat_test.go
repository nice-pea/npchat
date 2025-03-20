package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	baserepo "github.com/saime-0/nice-pea-chat/internal/repository/base"
)

func TestChats_Create(t *testing.T) {
	inputs := []ChatsCreateIn{
		{Name: "some_name", OwnerID: "qwe"},
	}
	for _, in := range inputs {
		// Настройка репозитория
		t.Run(fmt.Sprintf("%#v", in), func(t *testing.T) {
			chatsRepo := new(baserepo.ChatsRepository)
			membersRepo := new(baserepo.MembersRepository)
			chatsService := Chats{
				ChatsRepo:   chatsRepo,
				MembersRepo: membersRepo,
				History:     HistoryDummy{},
			}
			createdChat, err := chatsService.Create(in)
			assert.NoError(t, err)
			assert.NotEmpty(t, createdChat)
			assert.Equal(t, in.Name, createdChat.Name)
			assert.NotEmpty(t, createdChat.ID)
		})
	}
}
