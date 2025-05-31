package domain

import (
	"time"
)

// Чат

// Пользователь

type BasicAuth struct {
	Login    string // Логин пользователя
	Password string // Пароль пользователя
}

type OpenAuthLink struct {
	ExternalID string // ID пользователя провайдером
	Provider   string // Провайдер, которому принадлежит пользователь
	UserID     string
	Token      OpenAuthToken
}

type OpenAuthToken struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

type UserAggregate struct {
	ID   string // ID пользователя
	Name string // Имя пользователя
	Nick string // Ник пользователя

	BasicAuth     BasicAuth
	OpenAuthLinks OpenAuthLink
	//OpenAuthLinks map[string]OpenAuthLink
	Sessions []Session2
}

type Session2 struct {
	ID string // ID сессии
	//UserID string // ID пользователя, к которому относится сессия
	//Token  string // Токен сессии для аутентификации
	Status int // Статус сессии
}

// Сообщение
//
//type MessageAggregate struct {
//	ID      string // Уникальный идентификатор сообщения
//	Text    string // Текст сообщения
//	Timestamp time.Time // Время создания сообщения
//	UserID string // Идентификатор пользователя, отправившего сообщение
//	ChatID string // Идентификатор чата, к которому относится сообщение
//}
