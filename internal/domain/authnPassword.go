package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

type AuthnPassword struct {
	UserID   string
	Login    string
	Password string
}

const UserLoginMaxLen = 35
const UserLoginMinLen = 5

func (c AuthnPassword) ValidateLogin() error {
	// Проверка на длину логина
	if len([]rune(c.Login)) > UserLoginMaxLen {
		return fmt.Errorf("login cannot be longer than %d characters", UserLoginMaxLen)
	}
	if len([]rune(c.Login)) < UserLoginMinLen {
		return fmt.Errorf("login cannot be shorter than %d characters", UserLoginMinLen)
	}

	// Проверка на строку, состоящую только из пробелов
	if strings.TrimSpace(c.Login) == "" {
		return fmt.Errorf("login cannot consist of whitespace only")
	}

	// Проверка на первый символ
	if !isAllowedLastFirstLoginRune(rune(c.Login[0])) {
		return fmt.Errorf("login must start with a letter or digit")
	}
	// Проверка на последний символ
	if !isAllowedLastFirstLoginRune(rune(c.Login[len(c.Login)-1])) {
		return fmt.Errorf("login must trail with a letter or digit")
	}

	var hasLetters bool
	// Проверка каждого символа в имени пользователя
	for _, r := range c.Login {
		switch {
		case unicode.IsControl(r):
			return fmt.Errorf("login cannot contain control characters")
		case unicode.IsSpace(r):
			return fmt.Errorf("login cannot contain spaces")
		case !isAllowedLoginRune(r):
			return fmt.Errorf("login contains invalid characters")
		}
		if unicode.IsLetter(r) {
			hasLetters = true
		}
	}

	if !hasLetters {
		return fmt.Errorf("login must contain at least one letter or digit")
	}

	return nil
}

func isAllowedLoginRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isAllowedLastFirstLoginRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

const UserPasswordMaxLen = 128
const UserPasswordMinLen = 8

func (c AuthnPassword) ValidatePassword() error {
	// Проверка длины пароля
	if len(c.Password) < UserPasswordMinLen {
		return fmt.Errorf("password must be at least %d characters long", UserPasswordMinLen)
	}
	if len(c.Password) > UserPasswordMaxLen {
		return fmt.Errorf("password cannot exceed %d characters", UserPasswordMaxLen)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasInvalid bool
	)

	allowedSpecial := "~!?@#$%^&*_-+()[]{}></\\|\"'.,:;"

	for _, r := range c.Password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			// Проверка на арабские цифры (0-9)
			if r < '0' || r > '9' {
				return fmt.Errorf("only arabic digits (0-9) are allowed")
			}
			hasDigit = true
		case unicode.IsSpace(r):
			return fmt.Errorf("password cannot contain spaces")
		case strings.ContainsRune(allowedSpecial, r):
			// Разрешенные спецсимволы - ничего не делаем
		default:
			// Проверка на допустимые буквы (латиница + кириллица)
			if !isAllowedLetter(r) {
				hasInvalid = true
			}
		}
	}

	if hasInvalid {
		return fmt.Errorf("password contains invalid characters")
	}
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit (0-9)")
	}

	return nil
}

func isAllowedLetter(c rune) bool {
	// Разрешаем латинские и кириллические буквы
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
		(c >= 'а' && c <= 'я') || (c >= 'А' && c <= 'Я') ||
		c == 'ё' || c == 'Ё'
}

func (c AuthnPassword) ValidateUserID() error {
	if err := uuid.Validate(c.UserID); err != nil {
		return errors.Join(err, ErrInvitationUserIDValidate)
	}
	return nil
}

type AuthnPasswordRepository interface {
	Save(AuthnPassword) error
	List(filter AuthnPasswordFilter) ([]AuthnPassword, error)
}

type AuthnPasswordFilter struct {
	UserID   string
	Login    string
	Password string
}
