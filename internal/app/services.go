package app

import (
	acceptInvitation "github.com/nice-pea/npchat/internal/usecases/chats/accept_invitation"
	cancelInvitation "github.com/nice-pea/npchat/internal/usecases/chats/cancel_invitation"
	chatInvitations "github.com/nice-pea/npchat/internal/usecases/chats/chat_invitations"
	chatMembers "github.com/nice-pea/npchat/internal/usecases/chats/chat_members"
	createChat "github.com/nice-pea/npchat/internal/usecases/chats/create_chat"
	deleteMember "github.com/nice-pea/npchat/internal/usecases/chats/delete_member"
	leaveChat "github.com/nice-pea/npchat/internal/usecases/chats/leave_chat"
	myChats "github.com/nice-pea/npchat/internal/usecases/chats/my_chats"
	receivedInvitations "github.com/nice-pea/npchat/internal/usecases/chats/received_invitations"
	sendInvitation "github.com/nice-pea/npchat/internal/usecases/chats/send_invitation"
	updateName "github.com/nice-pea/npchat/internal/usecases/chats/update_name"
	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
	basicAuthLogin "github.com/nice-pea/npchat/internal/usecases/users/basic_auth/basic_auth_login"
	basicAuthRegistration "github.com/nice-pea/npchat/internal/usecases/users/basic_auth/basic_auth_registration"
	oauthAuthorize "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_authorize"
	oauthComplete "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_complete"
)

type usecasesBase struct {
	// Sessions

	*findSession.FindSessionsUsecase

	// Chats

	*acceptInvitation.AcceptInvitationUsecase
	*cancelInvitation.CancelInvitationUsecase
	*chatInvitations.ChatInvitationsUsecase
	*chatMembers.ChatMembersUsecase
	*createChat.CreateChatUsecase
	*deleteMember.DeleteMemberUsecase
	*leaveChat.LeaveChatUsecase
	*myChats.MyChatsUsecase
	*receivedInvitations.ReceivedInvitationsUsecase
	*sendInvitation.SendInvitationUsecase
	*updateName.UpdateNameUsecase

	// Users

	*basicAuthRegistration.BasicAuthRegistrationUsecase
	*basicAuthLogin.BasicAuthLoginUsecase
	*oauthAuthorize.OauthAuthorizeUsecase
	*oauthComplete.OauthCompleteUsecase
}

func initUsecases(rr *repositories, aa *adapters) usecasesBase {
	return usecasesBase{
		FindSessionsUsecase: &findSession.FindSessionsUsecase{
			Repo: rr.sessions,
		},
		AcceptInvitationUsecase: &acceptInvitation.AcceptInvitationUsecase{
			Repo:          rr.chats,
			EventConsumer: aa.eventBus,
		},
		CancelInvitationUsecase: &cancelInvitation.CancelInvitationUsecase{
			Repo:          rr.chats,
			EventConsumer: aa.eventBus,
		},
		ChatInvitationsUsecase: &chatInvitations.ChatInvitationsUsecase{
			Repo: rr.chats,
		},
		ChatMembersUsecase: &chatMembers.ChatMembersUsecase{
			Repo: rr.chats,
		},
		CreateChatUsecase: &createChat.CreateChatUsecase{
			Repo:          rr.chats,
			EventConsumer: aa.eventBus,
		},
		DeleteMemberUsecase: &deleteMember.DeleteMemberUsecase{
			Repo:          rr.chats,
			EventConsumer: aa.eventBus,
		},
		LeaveChatUsecase: &leaveChat.LeaveChatUsecase{
			Repo:          rr.chats,
			EventConsumer: aa.eventBus,
		},
		MyChatsUsecase: &myChats.MyChatsUsecase{
			Repo: rr.chats,
		},
		ReceivedInvitationsUsecase: &receivedInvitations.ReceivedInvitationsUsecase{
			Repo: rr.chats,
		},
		SendInvitationUsecase: &sendInvitation.SendInvitationUsecase{
			Repo:          rr.chats,
			EventConsumer: aa.eventBus,
		},
		UpdateNameUsecase: &updateName.UpdateNameUsecase{
			Repo:          rr.chats,
			EventConsumer: aa.eventBus,
		},
		BasicAuthRegistrationUsecase: &basicAuthRegistration.BasicAuthRegistrationUsecase{
			Repo:         rr.users,
			SessionsRepo: rr.sessions,
		},
		BasicAuthLoginUsecase: &basicAuthLogin.BasicAuthLoginUsecase{
			Repo:         rr.users,
			SessionsRepo: rr.sessions,
		},
		OauthAuthorizeUsecase: &oauthAuthorize.OauthAuthorizeUsecase{
			Providers: aa.oauthProviders,
		},
		OauthCompleteUsecase: &oauthComplete.OauthCompleteUsecase{
			Repo:         rr.users,
			Providers:    aa.oauthProviders,
			SessionsRepo: rr.sessions,
		},
	}
}
