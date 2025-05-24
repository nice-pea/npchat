package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

// AuthnPassword представляет собой метод аутентификации пользователя.
type AuthnPassword struct {
	UserID   string // Идентификатор пользователя
	Login    string // Логин пользователя
	Password string // Пароль пользователя
}

// UserLoginMaxLen максимальная длина логина пользователя.
const UserLoginMaxLen = 35

// UserLoginMinLen минимальная длина логина пользователя.
const UserLoginMinLen = 5

// ValidateLogin проверяет корректность логина пользователя.
func (c AuthnPassword) ValidateLogin() error {
	// Проверка на длину логина
	if len([]rune(c.Login)) > UserLoginMaxLen {
		return fmt.Errorf("логин не может быть длиннее %d символов", UserLoginMaxLen)
	}
	if len([]rune(c.Login)) < UserLoginMinLen {
		return fmt.Errorf("логин не может быть короче %d символов", UserLoginMinLen)
	}

	// Проверка на строку, состоящую только из пробелов
	if strings.TrimSpace(c.Login) == "" {
		return fmt.Errorf("логин не может состоять только из пробелов")
	}
	// Проверка на первый символ
	if !isAllowedLastFirstLoginRune(rune(c.Login[0])) {
		return fmt.Errorf("логин должен начинаться с буквы или цифры")
	}
	// Проверка на последний символ
	if !isAllowedLastFirstLoginRune(rune(c.Login[len(c.Login)-1])) {
		return fmt.Errorf("логин должен заканчиваться буквой или цифрой")
	}

	var hasLetters bool // Флаг для проверки наличия букв в логине
	// Проверка каждого символа в имени пользователя
	for _, r := range c.Login {
		switch {
		case unicode.IsControl(r):
			return fmt.Errorf("логин не может содержать управляющие символы")
		case unicode.IsSpace(r):
			return fmt.Errorf("логин не может содержать пробелы")
		case !isAllowedLoginRune(r):
			return fmt.Errorf("логин содержит недопустимые символы")
		}
		if unicode.IsLetter(r) {
			hasLetters = true // Устанавливаем флаг, если найдена буква
		}
	}

	if !hasLetters {
		return fmt.Errorf("логин должен содержать хотя бы одну букву или цифру")
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

// UserPasswordMaxLen максимальная длина пароля пользователя.
const UserPasswordMaxLen = 128

// UserPasswordMinLen минимальная длина пароля пользователя.
const UserPasswordMinLen = 8

// ValidatePassword проверяет корректность пароля пользователя.
func (c AuthnPassword) ValidatePassword() error {
	// Проверка длины пароля
	if len(c.Password) < UserPasswordMinLen {
		return fmt.Errorf("пароль должен быть не короче %d символов", UserPasswordMinLen)
	}
	if len(c.Password) > UserPasswordMaxLen {
		return fmt.Errorf("пароль не может быть длиннее %d символов", UserPasswordMaxLen)
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
	for _, r := range c.Password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true // Устанавливаем флаг, если найдена заглавная буква
		case unicode.IsLower(r):
			hasLower = true // Устанавливаем флаг, если найдена строчная буква
		case unicode.IsDigit(r):
			// Проверка на арабские цифры (0-9)
			if r < '0' || r > '9' {
				return fmt.Errorf("разрешены только арабские цифры (0-9)")
			}
			hasDigit = true // Устанавливаем флаг, если найдена цифра
		case unicode.IsSpace(r):
			return fmt.Errorf("пароль не может содержать пробелы") // Проверка на наличие пробелов
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
		return fmt.Errorf("пароль содержит недопустимые символы")
	}
	// Проверка на наличие заглавной буквы
	if !hasUpper {
		return fmt.Errorf("пароль должен содержать хотя бы одну заглавную букву")
	}
	// Проверка на наличие строчной буквы
	if !hasLower {
		return fmt.Errorf("пароль должен содержать хотя бы одну строчную букву")
	}
	// Проверка на наличие цифры
	if !hasDigit {
		return fmt.Errorf("пароль должен содержать хотя бы одну цифру (0-9)")
	}

	return nil // Пароль валиден
}

// isAllowedLetter проверяет, является ли символ допустимым (латиница или кириллица).
func isAllowedLetter(c rune) bool {
	// Разрешаем латинские и кириллические буквы
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= 'а' && c <= 'я') || (c >= 'А' && c <= 'Я') ||
		c == 'ё' || c == 'Ё'
}

// ValidateUserID проверяет корректность идентификатора пользователя.
func (c AuthnPassword) ValidateUserID() error {
	if err := uuid.Validate(c.UserID); err != nil {
		return errors.Join(err, ErrInvitationUserIDValidate) // Возвращаем ошибку, если идентификатор недействителен
	}
	return nil // Идентификатор валиден
}

// AuthnPasswordRepository интерфейс для работы с методом аутентификации по паролю.
type AuthnPasswordRepository interface {
	// Save сохраняет запись
	Save(AuthnPassword) error

	// List возвращает список с учетом фильтрации
	List(filter AuthnPasswordFilter) ([]AuthnPassword, error)
}

// AuthnPasswordFilter представляет собой фильтр по методам аутентификации по паролю.
type AuthnPasswordFilter struct {
	UserID   string // Идентификатор пользователя для фильтрации
	Login    string // Логин пользователя для фильтрации
	Password string // Пароль пользователя для фильтрации
}
