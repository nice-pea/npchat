package model

type Localization struct {
	Category string `json:"category"`
	Code     string `json:"code"`
	Locale   string `json:"locale"`
	Text     string `json:"text"`
}
