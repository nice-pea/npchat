package sessionn

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidateSessionName(t *testing.T) {
	tests := []struct {
		name        string
		sessionName string
		wantErr     bool
	}{
		{
			name:        "пустое значение",
			sessionName: "",
			wantErr:     true,
		},
		{
			name:        "любая строка",
			sessionName: "qwerty3456йцукенгш",
			wantErr:     false,
		},
		{
			name:        "браузер клиента",
			sessionName: gofakeit.UserAgent(),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSessionName(tt.name); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSessionStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		wantErr bool
	}{
		{
			name:    "пустой",
			status:  "",
			wantErr: true,
		},
		{
			name:    "какая-то строка",
			status:  "wertyu",
			wantErr: true,
		},
		{
			name:    "новая",
			status:  "new",
			wantErr: false,
		},
		{
			name:    "подтвержденная",
			status:  "verified",
			wantErr: false,
		},
		{
			name:    "истекшая",
			status:  "expired",
			wantErr: false,
		},
		{
			name:    "отозванная",
			status:  "revoked",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSessionStatus(tt.name); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
