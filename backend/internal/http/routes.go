package http

import (
	"fmt"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/app/optional"
	ucAuthn "github.com/saime-0/nice-pea-chat/internal/usecase/authn"
	ucLogin "github.com/saime-0/nice-pea-chat/internal/usecase/authn/login"
	ucChats "github.com/saime-0/nice-pea-chat/internal/usecase/chats"
	ucChatCreate "github.com/saime-0/nice-pea-chat/internal/usecase/chats/create"
	ucMembers "github.com/saime-0/nice-pea-chat/internal/usecase/members"
	ucMemberCreate "github.com/saime-0/nice-pea-chat/internal/usecase/members/create"
	ucMessages "github.com/saime-0/nice-pea-chat/internal/usecase/messages"
	ucMessageCreate "github.com/saime-0/nice-pea-chat/internal/usecase/messages/create"
	ucPermissions "github.com/saime-0/nice-pea-chat/internal/usecase/permissions"
	ucRoles "github.com/saime-0/nice-pea-chat/internal/usecase/roles"
	ucUsers "github.com/saime-0/nice-pea-chat/internal/usecase/users"
	ucUserCreate "github.com/saime-0/nice-pea-chat/internal/usecase/users/create"
	ucUserUpdate "github.com/saime-0/nice-pea-chat/internal/usecase/users/update"
)

func (s ServerParams) declareRoutes(muxHttp *http.ServeMux) {
	m := mux{ServeMux: muxHttp, s: s}
	// Service
	m.handle("/health", Health)
	// Users
	m.handle("/users", Users)
	m.handle("/users/update", UserUpdate)
	m.handle("/users/create", UserCreate)
	// Chats
	m.handle("/chats", Chats)
	m.handle("POST /chats/create", ChatCreate)
	// Messages
	m.handle("/messages", Messages)
	m.handle("/messages/create", MessageCreate)
	// Members
	m.handle("/members", Members)
	m.handle("/members/create", MemberCreate)
	m.handle("/permissions", Permissions)
	m.handle("/roles", Roles)
	// Authentication
	m.handle("/authn", Authn)
	m.handle("/authn/login", Login)
	m.handle("/credentials/", Login)
}

func MessageCreate(req Request) (any, error) {
	ucParams := ucMessageCreate.Params{
		DB: req.DB,
	}
	if err := parseJSONRequest(req.Body, &ucParams.Message); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func Messages(req Request) (_ any, err error) {
	ucParams := ucMessages.Params{
		IDs:        nil,
		ChatIDs:    nil,
		AuthorIDs:  nil,
		ReplyToIDs: nil,
		Boundary:   ucMessages.Boundary{},
		Limit:      optional.Uint{},
		DB:         req.DB,
	}
	if ucParams.IDs, err = uintsParam(req.Form, "ids"); err != nil {
		return nil, err
	}
	if ucParams.AuthorIDs, err = uintsParam(req.Form, "author_ids"); err != nil {
		return nil, err
	}
	if ucParams.ChatIDs, err = uintsParam(req.Form, "chat_ids"); err != nil {
		return nil, err
	}
	if ucParams.ReplyToIDs, err = uintsParam(req.Form, "reply_to_ids"); err != nil {
		return nil, err
	}

	if ucParams.Boundary.AroundID, err = uintOptionalParam(req.Form, "around_id"); err != nil {
		return nil, err
	} else if ucParams.Boundary.BeforeID, err = uintOptionalParam(req.Form, "before_id"); err != nil {
		return nil, err
	} else if ucParams.Boundary.AfterID, err = uintOptionalParam(req.Form, "after_id"); err != nil {
		return nil, err
	}

	if ucParams.Limit, err = uintOptionalParam(req.Form, "limit"); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func Members(req Request) (_ any, err error) {
	ucParams := ucMembers.Params{
		IDs:      nil,
		UserIDs:  nil,
		ChatIDs:  nil,
		IsPinned: 0,
		DB:       req.DB,
	}
	if ucParams.IDs, err = uintsParam(req.Form, "ids"); err != nil {
		return nil, err
	}
	if ucParams.UserIDs, err = uintsParam(req.Form, "user_ids"); err != nil {
		return nil, err
	}
	if ucParams.ChatIDs, err = uintsParam(req.Form, "chat_ids"); err != nil {
		return nil, err
	}
	if ucParams.IsPinned, err = boolOptionalParam(req.Form, "is_pinned"); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func MemberCreate(req Request) (any, error) {
	ucParams := ucMemberCreate.Params{
		DB: req.DB,
	}
	if err := parseJSONRequest(req.Body, &ucParams.Member); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func ChatCreate(req Request) (any, error) {
	ucParams := ucChatCreate.Params{
		DB: req.DB,
	}
	if err := parseJSONRequest(req.Body, &ucParams.Chat); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func Chats(req Request) (_ any, err error) {
	ucParams := ucChats.Params{
		IDs:                  nil,
		UserIDs:              nil,
		UnreadCounterForUser: optional.Uint{},
		Conn:                 req.PGXconn,
	}

	if ucParams.IDs, err = uintsParam(req.Form, "ids"); err != nil {
		return nil, err
	}
	if ucParams.UserIDs, err = uintsParam(req.Form, "user_ids"); err != nil {
		return nil, err
	}
	if ucParams.UnreadCounterForUser, err = uintOptionalParam(req.Form, "unread_counter"); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func UserUpdate(req Request) (any, error) {
	ucParams := ucUserUpdate.Params{
		DB: req.DB,
	}
	if err := parseJSONRequest(req.Body, &ucParams.User); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func UserCreate(req Request) (any, error) {
	ucParams := ucUserCreate.Params{
		DB: req.DB,
	}
	if err := parseJSONRequest(req.Body, &ucParams.User); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func Users(req Request) (any, error) {
	ucParams := ucUsers.Params{
		DB: req.DB,
	}

	return ucParams.Run()
}

func Login(req Request) (any, error) {
	ucParams := ucLogin.Params{
		Key: req.Form.Get("key"),
		DB:  req.DB,
	}

	return ucParams.Run()
}

func Authn(req Request) (any, error) {
	ucParams := ucAuthn.Params{
		Token: req.Token,
		DB:    req.DB,
	}

	return ucParams.Run()
}

func Permissions(req Request) (any, error) {
	ucParams := ucPermissions.Params{
		Locale: req.Locale,
		L10n:   req.L10n,
	}

	return ucParams.Run()
}

func Roles(req Request) (_ any, err error) {
	// Load params
	ucParams := ucRoles.Params{
		IDs:  nil,
		Name: req.Form.Get("name"),
		DB:   req.DB,
	}
	if ucParams.IDs, err = uintsParam(req.Form, "ids"); err != nil {
		return nil, err
	}

	return ucParams.Run()
}

func Health(req Request) (any, error) {
	return req.L10n.Localize("none:ok", req.Locale, nil)
}

func userAuthn(req Request) (any, error) {
	users, err := ucUsers.Params{DB: req.DB}.Run()
	if err != nil {
		return nil, err
	} else if len(users) != 1 {
		return nil, fmt.Errorf("expected 1 user, got %d", len(users))
	}

	return users[0], nil
}
