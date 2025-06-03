package domain

import (
	"time"
)

// Сессия

type Session2 struct {
	ID           string // ID сессии
	UserID       string
	Name         string
	Status       int // Статус сессии
	AccessToken  SessionToken
	RefreshToken SessionToken

	//Client SessionClient
}

type SessionToken struct {
	Token  string
	Expiry time.Time
}

//type SessionClient struct {
//	ID   string // Идентификатор устройства
//	Type string // android/ios/web
//}

// Push-уведомления
//
//type PushConfig struct {
//	UserID   string
//	Channels []PushChannel // Имена каналов, на которые отправляются уведомления
//}
//
//type PushChannel struct {
//	Name     string       // Имя канала
//	Provider PushProvider // ZeroValue - для глобального значения
//	Enabled  bool         // Включено ли уведомление в канале
//}
//
//type PushProvider struct {
//	Name  string // apn, fcm, gcm, sms
//	Token string // Токен для отправки уведомлений
//}

// Сообщение
//
//type MessageAggregate struct {
//	ID      string // Уникальный идентификатор сообщения
//	Text    string // Текст сообщения
//	Timestamp time.Time // Время создания сообщения
//	UserID string // Идентификатор пользователя, отправившего сообщение
//	ChatID string // Идентификатор чата, к которому относится сообщение
//}
