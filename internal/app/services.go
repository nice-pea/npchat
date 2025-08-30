package app

import (
	acceptInvitation "github.com/nice-pea/npchat/internal/service/chats/accept_invitation"
	cancelInvitation "github.com/nice-pea/npchat/internal/service/chats/cancel_invitation"
	chatInvitations "github.com/nice-pea/npchat/internal/service/chats/chat_invitations"
	chatMembers "github.com/nice-pea/npchat/internal/service/chats/chat_members"
	createChat "github.com/nice-pea/npchat/internal/service/chats/create_chat"
	deleteMember "github.com/nice-pea/npchat/internal/service/chats/delete_member"
	leaveChat "github.com/nice-pea/npchat/internal/service/chats/leave_chat"
	myChats "github.com/nice-pea/npchat/internal/service/chats/my_chats"
	receivedInvitations "github.com/nice-pea/npchat/internal/service/chats/received_invitations"
	sendInvitation "github.com/nice-pea/npchat/internal/service/chats/send_invitation"
	updateName "github.com/nice-pea/npchat/internal/service/chats/update_name"
	findSession "github.com/nice-pea/npchat/internal/service/sessions/find_session"
	basicAuthLogin "github.com/nice-pea/npchat/internal/service/users/basic_auth/basic_auth_login"
	basicAuthRegistration "github.com/nice-pea/npchat/internal/service/users/basic_auth/basic_auth_registration"
	"github.com/nice-pea/npchat/internal/service/users/oauth/completeOAuthRegistration"
	completeOAuthLogin "github.com/nice-pea/npchat/internal/service/users/oauth/complete_oauth_login"
	initOAuthLogin "github.com/nice-pea/npchat/internal/service/users/oauth/init_oauth_login"
	initOAuthRegistration "github.com/nice-pea/npchat/internal/service/users/oauth/init_oauth_registration"
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
	*initOAuthRegistration.InitOAuthRegistrationUsecase
	*completeOAuthRegistration.CompleteOAuthRegistrationUsecase
	*initOAuthLogin.InitOAuthLoginUsecase
	*completeOAuthLogin.CompleteOAuthLoginUsecase
}

func initUsecases(rr *repositories, aa *adapters) usecasesBase {
	return usecasesBase{
		FindSessionsUsecase: &findSession.FindSessionsUsecase{
			Repo: rr.sessions,
		},
		AcceptInvitationUsecase: &acceptInvitation.AcceptInvitationUsecase{
			Repo: rr.chats,
		},
		CancelInvitationUsecase: &cancelInvitation.CancelInvitationUsecase{
			Repo: rr.chats,
		},
		ChatInvitationsUsecase: &chatInvitations.ChatInvitationsUsecase{
			Repo: rr.chats,
		},
		ChatMembersUsecase: &chatMembers.ChatMembersUsecase{
			Repo: rr.chats,
		},
		CreateChatUsecase: &createChat.CreateChatUsecase{
			Repo: rr.chats,
		},
		DeleteMemberUsecase: &deleteMember.DeleteMemberUsecase{
			Repo: rr.chats,
		},
		LeaveChatUsecase: &leaveChat.LeaveChatUsecase{
			Repo: rr.chats,
		},
		MyChatsUsecase: &myChats.MyChatsUsecase{
			Repo: rr.chats,
		},
		ReceivedInvitationsUsecase: &receivedInvitations.ReceivedInvitationsUsecase{
			Repo: rr.chats,
		},
		SendInvitationUsecase: &sendInvitation.SendInvitationUsecase{
			Repo: rr.chats,
		},
		UpdateNameUsecase: &updateName.UpdateNameUsecase{
			Repo: rr.chats,
		},
		BasicAuthRegistrationUsecase: &basicAuthRegistration.BasicAuthRegistrationUsecase{
			Repo:         rr.users,
			SessionsRepo: rr.sessions,
		},
		BasicAuthLoginUsecase: &basicAuthLogin.BasicAuthLoginUsecase{
			Repo:         rr.users,
			SessionsRepo: rr.sessions,
		},
		InitOAuthRegistrationUsecase: &initOAuthRegistration.InitOAuthRegistrationUsecase{
			Providers: aa.oauthProviders,
		},
		CompleteOAuthRegistrationUsecase: &completeOAuthRegistration.CompleteOAuthRegistrationUsecase{
			Repo:         rr.users,
			Providers:    aa.oauthProviders,
			SessionsRepo: rr.sessions,
		},
		InitOAuthLoginUsecase:     &initOAuthLogin.InitOAuthLoginUsecase{},
		CompleteOAuthLoginUsecase: &completeOAuthLogin.CompleteOAuthLoginUsecase{},
	}
}
