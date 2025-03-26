package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestInvitation_ValidateId(t *testing.T) {
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
			name:    "пустая строка как id",
			fields:  fields{ID: "", ChatID: ""},
			wantErr: true,
		},
		{
			name:    "коротка строка точно не uuid",
			fields:  fields{ID: "fndsef", ChatID: ""},
			wantErr: true,
		},
		{
			name:    "короткая строка из символов точно не uuid",
			fields:  fields{ID: "----", ChatID: ""},
			wantErr: true,
		},
		{
			name:    "строка из символов точно не uuid",
			fields:  fields{ID: ",,,,,,,,", ChatID: ""},
			wantErr: true,
		},
		{
			name:    "нужное количество символов",
			fields:  fields{ID: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", ChatID: ""},
			wantErr: true,
		},
		{
			name:    "это uuid",
			fields:  fields{ID: "1cee9c74-a359-425c-b1bb-91c8a35e7b21", ChatID: ""},
			wantErr: false,
		},
		{
			name:    "это uuid",
			fields:  fields{ID: "0195ba16-f44c-7a3b-b326-94697ec6b00e", ChatID: ""},
			wantErr: false,
		},
		{
			name:    "uuid генерируемый библиотекой",
			fields:  fields{ID: uuid.NewString(), ChatID: ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Invitation {
				ID:     tt.fields.ID,
				ChatID: tt.fields.ChatID,
			}
			if err := i.ValidateID(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
