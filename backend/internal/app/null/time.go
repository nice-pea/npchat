package null

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Time struct {
	sql.NullTime
}

func (t Time) IsZero() bool {
	return !t.Valid
}

// MarshalJSON - пользовательская реализация маршалинга
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return json.Marshal(nil) // Если значение недействительно, возвращаем null
	}
	return json.Marshal(t.Time) // Иначе маршалим время
}

// UnmarshalJSON - пользовательская реализация анмаршалинга
func (t *Time) UnmarshalJSON(data []byte) error {
	var aux *time.Time // Временная переменная для хранения значения

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux == nil {
		t.Valid = false // Если значение null, устанавливаем Valid в false
		return nil
	}

	t.Time = *aux
	t.Valid = true // Устанавливаем Valid в true, если значение не null
	return nil
}
