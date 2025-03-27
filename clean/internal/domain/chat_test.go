package domain

import (
	"strings"
	"testing"
)

func TestChat_ValidateID(t *testing.T) {
	RunValidateRequiredIDTest(t, func(ID string) error {
		c := Chat{ID: ID}
		return c.ValidateID()
	})
}

func TestChat_ValidateName(t *testing.T) {
	type fields struct {
		ID   string
		Name string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "пустая строка",
			fields:  fields{Name: ""},
			wantErr: true,
		},
		{
			name:    "превышает лимит в 50 символов",
			fields:  fields{Name: strings.Repeat("a", 51)},
			wantErr: true,
		},
		{
			name:    "содержит пробел в начале",
			fields:  fields{Name: " name"},
			wantErr: true,
		},
		{
			name:    "содержит пробел в конце",
			fields:  fields{Name: "name "},
			wantErr: true,
		},
		{
			name:    "содержит таб",
			fields:  fields{Name: "na\tme"},
			wantErr: true,
		},
		{
			name:    "содержит новую строку",
			fields:  fields{Name: "na\nme"},
			wantErr: true,
		},
		{
			name:    "содержит цифры",
			fields:  fields{Name: "1na13me4"},
			wantErr: false,
		},
		{
			name:    "содержит пробел в середине",
			fields:  fields{Name: "na me"},
			wantErr: false,
		},
		{
			name:    "содержит пробелы в середине",
			fields:  fields{Name: "na  me"},
			wantErr: false,
		},
		{
			name:    "содержит знаки",
			fields:  fields{Name: "??na??me.#1432&^$(@"},
			wantErr: false,
		},
		{
			name:    "содержит только знаки",
			fields:  fields{Name: "?>><#(*@$&"},
			wantErr: false,
		},
		{
			name:    "содержит только пробелы",
			fields:  fields{Name: " "},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Chat{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			if err := c.ValidateName(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChat_ValidateChiefUserID(t *testing.T) {
	RunValidateRequiredIDTest(t, func(ChiefUserID string) error {
		c := Chat{
			ChiefUserID: ChiefUserID,
		}
		return c.ValidateChiefUserID()
	})
}
