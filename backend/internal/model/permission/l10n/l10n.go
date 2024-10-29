package l10n

import (
	"fmt"
)

const Category = "permissions"

func Name(permission uint8) string {
	return fmt.Sprintf("%s:name%d", Category, permission)
}

func Desc(permission uint8) string {
	return fmt.Sprintf("%s:desc%d", Category, permission)
}
