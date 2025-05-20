package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AuthnPassword_ValidateLogin(t *testing.T) {
	tests := []struct {
		testname string
		login    string
		wantErr  bool
	}{
		// Invalid cases
		{testname: "empty login", login: "", wantErr: true},
		{testname: "whitespace only", login: " ", wantErr: true},
		{testname: "control characters", login: "\n\t\r\a\f\v", wantErr: true},
		{testname: "space in middle", login: "first last", wantErr: true},
		{testname: "surrounding spaces", login: " login ", wantErr: true},
		{testname: "digits only", login: "11111", wantErr: true},
		{testname: "special characters", login: "-!@#$%^&*()+}{:\"'?><.", wantErr: true},
		{testname: "trailing underscore", login: "login_", wantErr: true},
		{testname: "leading underscore", login: "_login", wantErr: true},
		{testname: "hyphen in middle", login: "first-last", wantErr: true},
		{testname: "cyrillic letters", login: "Логин", wantErr: true},
		{testname: "chinese characters", login: "名字名字名字", wantErr: true},
		{testname: "emoji", login: "😊nick", wantErr: true}, // или true, если эмодзи запрещены
		{testname: "special unicode (ñ, é)", login: "niño café", wantErr: true},
		{testname: "japanese (kanji + hiragana)", login: "名前なまえ", wantErr: true},
		{testname: "too long login", login: strings.Repeat("a", 35+1), wantErr: true},
		{testname: "too short login", login: strings.Repeat("a", 5-1), wantErr: true},

		// Valid cases
		{testname: "valid simple login", login: "login", wantErr: false},
		{testname: "valid login with digits", login: "1lo1gin1", wantErr: false},
		{testname: "valid underscore separated", login: "login_login", wantErr: false},
		{testname: "alphabetic", login: "abcdefghijklmnopqrstuvwxyz", wantErr: false},
		{testname: "alphabetic", login: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			ap := &AuthnPassword{Login: tt.login}
			if err := ap.ValidateLogin(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_AuthnPassword_ValidatePassword(t *testing.T) {
	tests := []struct {
		testname string
		password string
		wantErr  bool
	}{
		// Невалидные пароли
		{testname: "empty password", password: "", wantErr: true},
		{testname: "too short (7 chars)", password: "Ab1!xyz", wantErr: true},
		{testname: "too long (129 chars)", password: strings.Repeat("A", 128+1), wantErr: true},
		{testname: "no uppercase", password: "abc123!@#", wantErr: true},
		{testname: "no lowercase", password: "ABC123!@#", wantErr: true},
		{testname: "no digits", password: "Abcdef!@#", wantErr: true},
		{testname: "contains spaces", password: "Abc 123!@#", wantErr: true},
		{testname: "non-arabic digits", password: "Abc١٢٣!@#", wantErr: true}, // арабские цифры ١٢٣
		{testname: "invalid symbols", password: "Abc123€§¶", wantErr: true},

		// Валидные пароли
		{testname: "min length (8 chars)", password: "Ab1!xyzZ", wantErr: false},
		{testname: "with latin letters", password: "Password123!", wantErr: false},
		{testname: "with cyrillic letters", password: "Пароль123!", wantErr: false},
		{testname: "all allowed special chars", password: "Aa1~!?@#$%^&*_-+()[]{}></\\|\"'.,:;", wantErr: false},
		{testname: "complex password", password: "P@ssw0rd_123", wantErr: false},
		{testname: "mixed languages", password: "Passворд123!", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			ap := &AuthnPassword{Password: tt.password}
			if err := ap.ValidatePassword(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_AuthnPassword_ValidateUserID(t *testing.T) {
	tests := []struct {
		testname string
		userID   string
		wantErr  bool
	}{
		{
			testname: "empty user ID",
			userID:   "",
			wantErr:  true,
		},
		{
			testname: "invalid UUID format",
			userID:   "not-a-uuid",
			wantErr:  true,
		},
		{
			testname: "UUID with invalid characters",
			userID:   "123e4567-e89b-12d3-a456-42661417400g",
			wantErr:  true,
		},
		{
			testname: "UUID wrong length",
			userID:   "123e4567-e89b-12d3-a456",
			wantErr:  true,
		},
		{
			testname: "UUID with missing hyphens",
			userID:   "123e4567e89b12d3a456426614174000",
			wantErr:  false,
		},
		{
			testname: "valid user ID",
			userID:   "123e4567-e89b-12d3-a456-426614174000",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			ap := &AuthnPassword{UserID: tt.userID}
			if err := ap.ValidateUserID(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
