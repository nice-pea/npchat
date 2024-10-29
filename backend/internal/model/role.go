package model

type Role struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"-" json:"name"`
	ChatID      uint    `json:"chat_id"`
	Permissions []uint8 `gorm:"type:smallint[];default:'{}'" json:"permissions,omitempty"`
}
