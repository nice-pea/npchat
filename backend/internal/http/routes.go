package http

import (
	"net/http"

	ucAuthn "github.com/saime-0/nice-pea-chat/internal/usecase/authn"
	ucLogin "github.com/saime-0/nice-pea-chat/internal/usecase/authn/login"
	ucChats "github.com/saime-0/nice-pea-chat/internal/usecase/chats"
	ucChatCreate "github.com/saime-0/nice-pea-chat/internal/usecase/chats/create"
	ucPermissions "github.com/saime-0/nice-pea-chat/internal/usecase/permissions"
	ucRoles "github.com/saime-0/nice-pea-chat/internal/usecase/roles"
	"github.com/saime-0/nice-pea-chat/internal/usecase/users"
	ucUserUpdate "github.com/saime-0/nice-pea-chat/internal/usecase/users/update"
)

func (s ServerParams) declareRoutes(muxHttp *http.ServeMux) {
	m := mux{ServeMux: muxHttp, s: s}
	// Service
	m.handle("/health", Health)
	// Users
	m.handle("/users", Users)
	m.handle("/users/update", UserUpdate)
	// Chats
	m.handle("/chats", Chats)
	m.handle("POST /chats/create", ChatCreate)
	m.handle("/permissions", Permissions)
	m.handle("/roles", Roles)
	// Authentication
	m.handle("/authn", Authn)
	m.handle("/authn/login", Login)
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
		IDs:     nil,
		UserIDs: nil,
		DB:      req.DB,
	}

	if ucParams.IDs, err = uintsParam(req.Form, "ids"); err != nil {
		return nil, err
	}
	if ucParams.UserIDs, err = uintsParam(req.Form, "user_ids"); err != nil {
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

func Users(req Request) (any, error) {
	ucParams := users.Params{
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
		Token: req.Form.Get("token"),
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
