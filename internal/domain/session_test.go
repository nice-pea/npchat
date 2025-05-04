package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

func Test_Session_ValidateID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(ID string) error {
		s := Session{ID: ID}
		return s.ValidateID()
	})
}
func Test_Session_ValidateUserID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(ID string) error {
		s := Session{UserID: ID}
		return s.ValidateUserID()
	})
}
func Test_Session_ValidateToken(t *testing.T) {}
func Test_Session_ValidateStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{
			name:    "допустимый статус - новый",
			status:  SessionStatusNew,
			wantErr: false,
		},
		{
			name:    "допустимый статус - в ожидании",
			status:  SessionStatusPending,
			wantErr: false,
		},
		{
			name:    "допустимый статус - проверенный",
			status:  SessionStatusVerified,
			wantErr: false,
		},
		{
			name:    "допустимый статус - просроченный",
			status:  SessionStatusExpired,
			wantErr: false,
		},
		{
			name:    "допустимый статус - отозванный",
			status:  SessionStatusRevoked,
			wantErr: false,
		},
		{
			name:    "допустимый статус - неудачный",
			status:  SessionStatusFailed,
			wantErr: false,
		},
		{
			name:    "недопустимый статус - нулевой",
			status:  0,
			wantErr: true,
		},
		{
			name:    "недопустимый статус - отрицательный",
			status:  -1,
			wantErr: true,
		},
		{
			name:    "недопустимый статус - слишком большой",
			status:  7,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Session{Status: tt.status}
			if err := s.ValidateStatus(); tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
