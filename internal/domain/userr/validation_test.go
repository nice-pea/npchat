package userr

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateBasicAuthLogin(t *testing.T) {
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
			if err := ValidateBasicAuthLogin(tt.login); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBasicAuthPassword(t *testing.T) {
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
			if err := ValidateBasicAuthPassword(tt.password); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUserName(t *testing.T) {
	tests := []struct {
		testname string
		name     string
		wantErr  bool
	}{
		// Invalid cases
		{testname: "empty name", name: "", wantErr: true},
		{testname: "whitespace only", name: " ", wantErr: true},
		{testname: "control characters", name: "q\n\t\r\a\f\vq", wantErr: true},
		{testname: "leading space", name: " Name", wantErr: true},
		{testname: "trailing space", name: "Name ", wantErr: true},
		{testname: "too long name", name: strings.Repeat("名", 35+1), wantErr: true},

		// Valid ASCII cases
		{testname: "single character", name: "q", wantErr: false},
		{testname: "single digit", name: "1", wantErr: false},
		{testname: "lowercase name", name: "name", wantErr: false},
		{testname: "alphanumeric name", name: "Name2", wantErr: false},
		{testname: "name with space", name: "first last", wantErr: false},

		// Unicode cases (valid if supported)
		{testname: "cyrillic letters", name: "ИмяФамилия", wantErr: false},
		{testname: "chinese characters", name: "名字", wantErr: false},
		{testname: "emoji", name: "😊nick", wantErr: false},
		{testname: "special unicode (ñ, é)", name: "niño café", wantErr: false},
		{testname: "japanese (kanji + hiragana)", name: "名前なまえ", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			if err := ValidateUserName(tt.name); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUserNick(t *testing.T) {
	tests := []struct {
		testname string
		nick     string
		wantErr  bool
	}{
		// Invalid cases
		{testname: "whitespace only", nick: " ", wantErr: true},
		{testname: "control characters", nick: "\n\t\r\a\f\v", wantErr: true},
		{testname: "space in middle", nick: "first last", wantErr: true},
		{testname: "surrounding spaces", nick: " name ", wantErr: true},
		{testname: "digits only", nick: "1111", wantErr: true},
		{testname: "special characters", nick: "-!@#$%^&*()+}{:\"'?><.", wantErr: true},
		{testname: "trailing underscore", nick: "name_", wantErr: true},
		{testname: "leading underscore", nick: "_name", wantErr: true},
		{testname: "hyphen in middle", nick: "first-last", wantErr: true},
		{testname: "cyrillic letters", nick: "ИмяФамилия", wantErr: true},
		{testname: "chinese characters", nick: "名字", wantErr: true},
		{testname: "emoji", nick: "😊nick", wantErr: true}, // или true, если эмодзи запрещены
		{testname: "special unicode (ñ, é)", nick: "niño café", wantErr: true},
		{testname: "japanese (kanji + hiragana)", nick: "名前なまえ", wantErr: true},
		{testname: "too long name", nick: strings.Repeat("a", 35+1), wantErr: true},

		// Valid cases
		{testname: "empty nick", nick: "", wantErr: false},
		{testname: "valid simple name", nick: "name", wantErr: false},
		{testname: "valid name with digits", nick: "1na1me1", wantErr: false},
		{testname: "valid underscore separated", nick: "first_last", wantErr: false},
		{testname: "alphabetic", nick: "abcdefghijklmnopqrstuvwxyz", wantErr: false},
		{testname: "alphabetic", nick: "ABCDEFGHIJKLMNOPQRSTUVWXYZ", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			if err := ValidateUserNick(tt.nick); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
