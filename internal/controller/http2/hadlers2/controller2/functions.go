package controller2

import (
	"encoding/json"
	"fmt"
)

func DecodeBody[T any](context Context, dst *T) error {
	if err := json.NewDecoder(context.Request().Body).Decode(dst); err != nil {
		return fmt.Errorf("json Decoder Decode: %w", err)
	}

	return nil
}
