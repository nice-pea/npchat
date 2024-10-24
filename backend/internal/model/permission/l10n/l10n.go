package l10n

import "strconv"

const Category = "permissions"

func Code(permission int) string {
	return Category + ":" + strconv.Itoa(permission)
}
