package userr

import (
	"strings"
	"unicode"
)

// ValidateBasicAuthPassword проверяет корректность пароля пользователя.
func ValidateBasicAuthPassword(password string) error {
	// Проверка длины пароля
	if len(password) < UserPasswordMinLen {
		return ErrPasswordTooShort
	}
	if len(password) > UserPasswordMaxLen {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   bool // Флаг наличия заглавной буквы
		hasLower   bool // Флаг наличия строчной буквы
		hasDigit   bool // Флаг наличия цифры
		hasInvalid bool // Флаг наличия недопустимых символов
	)

	// Разрешенные специальные символы
	allowedSpecial := "~!?@#$%^&*_-+()[]{}></\\|\"'.,:;"

	// Проверка каждого символа в пароле
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true // Устанавливаем флаг, если найдена заглавная буква
		case unicode.IsLower(r):
			hasLower = true // Устанавливаем флаг, если найдена строчная буква
		case unicode.IsDigit(r):
			// Проверка на арабские цифры (0-9)
			if r < '0' || r > '9' {
				return ErrOnlyArabicDigits
			}
			hasDigit = true // Устанавливаем флаг, если найдена цифра
		case unicode.IsSpace(r):
			return ErrPasswordContainsSpaces // Проверка на наличие пробелов
		case strings.ContainsRune(allowedSpecial, r):
			// Разрешенные спецсимволы - ничего не делаем
		default:
			// Проверка на допустимые буквы (латиница + кириллица)
			if !isAllowedLetter(r) {
				hasInvalid = true // Устанавливаем флаг, если найден недопустимый символ
			}
		}
	}

	// Проверка на наличие недопустимых символов
	if hasInvalid {
		return ErrPasswordInvalidChars
	}
	// Проверка на наличие заглавной буквы
	if !hasUpper {
		return ErrPasswordNoUppercase
	}
	// Проверка на наличие строчной буквы
	if !hasLower {
		return ErrPasswordNoLowercase
	}
	// Проверка на наличие цифры
	if !hasDigit {
		return ErrPasswordNoDigit
	}

	return nil // Пароль валиден
}

// UserPasswordMaxLen максимальная длина пароля пользователя.
const UserPasswordMaxLen = 128

// UserPasswordMinLen минимальная длина пароля пользователя.
const UserPasswordMinLen = 8

// isAllowedLetter проверяет, является ли символ допустимым (латиница или кириллица).
func isAllowedLetter(c rune) bool {
	// Разрешаем латинские и кириллические буквы
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= 'а' && c <= 'я') || (c >= 'А' && c <= 'Я') ||
		c == 'ё' || c == 'Ё'
}

// UserLoginMaxLen максимальная длина логина пользователя.
const UserLoginMaxLen = 35

// UserLoginMinLen минимальная длина логина пользователя.
const UserLoginMinLen = 5

// ValidateBasicAuthLogin проверяет корректность логина пользователя.
func ValidateBasicAuthLogin(login string) error {
	// Проверка на длину логина
	if len([]rune(login)) > UserLoginMaxLen {
		return ErrLoginTooLong
	}
	if len([]rune(login)) < UserLoginMinLen {
		return ErrLoginTooShort
	}

	// Проверка на строку, состоящую только из пробелов
	if strings.TrimSpace(login) == "" {
		return ErrLoginOnlySpaces
	}
	// Проверка на первый символ
	if !isAllowedLastFirstLoginRune(rune(login[0])) {
		return ErrLoginStartChar
	}
	// Проверка на последний символ
	if !isAllowedLastFirstLoginRune(rune(login[len(login)-1])) {
		return ErrLoginEndChar
	}

	var hasLetters bool // Флаг для проверки наличия букв в логине
	// Проверка каждого символа в имени пользователя
	for _, r := range login {
		switch {
		case unicode.IsControl(r):
			return ErrLoginControlChars
		case unicode.IsSpace(r):
			return ErrLoginSpaces
		case !isAllowedLoginRune(r):
			return ErrLoginInvalidChars
		}
		if unicode.IsLetter(r) {
			hasLetters = true // Устанавливаем флаг, если найдена буква
		}
	}

	if !hasLetters {
		return ErrLoginNoLetters
	}

	return nil // Логин валиден
}

// isAllowedLoginRune проверяет, является ли символ допустимым для логина.
func isAllowedLoginRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isAllowedLastFirstLoginRune проверяет, является ли символ допустимым для первого или последнего символа логина.
func isAllowedLastFirstLoginRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// UserNameMaxLen максимальная длина имени пользователя.
const UserNameMaxLen = 35

// ValidateUserName проверяет корректность имени пользователя.
func ValidateUserName(name string) error {
	// Проверить, не является ли имя пустым или содержит только пробелы
	if strings.TrimSpace(name) == "" {
		return ErrNameEmpty
	}

	// Проверка на длину имени пользователя
	if len([]rune(name)) > UserNameMaxLen {
		return ErrNameTooLong
	}

	// Check for leading or trailing spaces
	if name != strings.TrimSpace(name) {
		return ErrNameSpaces
	}

	// Проверка на управляющие символы (кроме обычного пробела)
	for _, r := range name {
		if unicode.IsControl(r) && r != ' ' { // \n, \t, \r и т.д.
			return ErrNameControlChars
		}
	}

	return nil
}

// UserNickMaxLen максимальная длина ника пользователя.
const UserNickMaxLen = 35

// ValidateUserNick проверяет корректность ника пользователя.
func ValidateUserNick(nick string) error {
	// Проверка на пустое имя пользователя
	if nick == "" {
		return nil
	}

	// Проверка на длину имени пользователя
	if len([]rune(nick)) > UserNickMaxLen {
		return ErrNickTooLong
	}

	// Проверка на строку, состоящую только из пробелов
	if strings.TrimSpace(nick) == "" {
		return ErrNickOnlySpaces
	}

	// Проверка на первый символ
	if !isAllowedLastFirstNickRune(rune(nick[0])) {
		return ErrNickStartChar
	}
	// Проверка на последний символ
	if !isAllowedLastFirstNickRune(rune(nick[len(nick)-1])) {
		return ErrNickEndChar
	}

	var hasLetters bool
	// Проверка каждого символа в имени пользователя
	for _, r := range nick {
		switch {
		case unicode.IsControl(r):
			return ErrNickControlChars
		case unicode.IsSpace(r):
			return ErrNickSpaces
		case !isAllowedNickRune(r):
			return ErrNickInvalidChars
		}
		if unicode.IsLetter(r) {
			hasLetters = true
		}
	}

	if !hasLetters {
		return ErrNickNoLetters
	}

	return nil // Ник валиден
}

// isAllowedNickRune проверяет, разрешен ли символ в имени пользователя
func isAllowedNickRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// isAllowedLastFirstNickRune проверяет, разрешен ли символ в качестве первого или последнего символа ника
func isAllowedLastFirstNickRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
