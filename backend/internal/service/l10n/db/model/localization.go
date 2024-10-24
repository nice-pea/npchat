package model

type Localization struct {
	Category string `json:"category"`
	Item     string `json:"item"`
	Locale   string `json:"locale"`
	Text     string `json:"text"`
}
