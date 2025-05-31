package chatt

import (
	"regexp"
)

// ValidateChatName проверяет корректность названия чата.
func ValidateChatName(name string) error {
	// Регулярное выражение для проверки названия чата
	var chatNameRegexp = regexp.MustCompile(`^[^\s\n\t][^\n\t]{0,48}[^\s\n\t]$`)
	if !chatNameRegexp.MatchString(name) {
		return ErrChatNameValidate // Возвращает ошибку, если название некорректно
	}

	return nil // Название валидно
}
