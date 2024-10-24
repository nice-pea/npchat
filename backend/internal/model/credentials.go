package model

type Credentials struct {
	UserID uint64 `gorm:"primaryKey" json:"user_id"`
	Key    string `json:"key"`
}
