package l10n

import (
	"fmt"
)

const Category = "permissions"

func Name(permission int) string {
	return fmt.Sprintf("%s:name%d", Category, permission)
}

func Desc(permission int) string {
	return fmt.Sprintf("%s:desc%d", Category, permission)
}
