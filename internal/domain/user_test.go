package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

func Test_User_ValidateID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(ID string) error {
		u := User{ID: ID}
		return u.ValidateID()
	})
}

func Test_User_ValidateNick(t *testing.T) {
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
			u := User{Nick: tt.nick}
			if err := u.ValidateNick(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_UserValidateName(t *testing.T) {
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
			u := User{Name: tt.name}
			if err := u.ValidateName(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
