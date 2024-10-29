package model

type Credentials struct {
	UserID uint   `gorm:"primaryKey" json:"user_id"`
	Key    string `json:"key"`
}
