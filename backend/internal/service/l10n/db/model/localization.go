package model

type Localization struct {
	Category string `gorm:"" json:"category"`
	Code     string `gorm:"" json:"code"`
	Locale   string `gorm:"" json:"locale"`
	Text     string `gorm:"" json:"text"`
}
