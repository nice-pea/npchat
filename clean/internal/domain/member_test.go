package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestMember_ValidateID(t *testing.T) {
	RunValidateIDTest(t, func(ID string) error {
		m := Member{ID: ID}
		return m.ValidateID()
	})
}

func TestMember_ValidateChatID(t *testing.T) {
	type fields struct {
		ID     string
		ChatID string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "пустая строка",
			fields:  fields{ChatID: ""},
			wantErr: true,
		},
		{
			name:    "коротка строка точно не uuid",
			fields:  fields{ChatID: "fndsef"},
			wantErr: true,
		},
		{
			name:    "короткая строка из символов точно не uuid",
			fields:  fields{ChatID: "----"},
			wantErr: true,
		},
		{
			name:    "строка из символов точно не uuid",
			fields:  fields{ChatID: ",,,,,,,,"},
			wantErr: true,
		},
		{
			name:    "нужное количество символов",
			fields:  fields{ChatID: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
			wantErr: true,
		},
		{
			name:    "это uuid",
			fields:  fields{ChatID: "1cee9c74-a359-425c-b1bb-91c8a35e7b21"},
			wantErr: false,
		},
		{
			name:    "это uuid",
			fields:  fields{ChatID: "0195ba16-f44c-7a3b-b326-94697ec6b00e"},
			wantErr: false,
		},
		{
			name:    "uuid генерируемый библиотекой",
			fields:  fields{ChatID: uuid.NewString()},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Member{
				ID:     tt.fields.ID,
				ChatID: tt.fields.ChatID,
			}
			if err := m.ValidateChatID(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateChatID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
