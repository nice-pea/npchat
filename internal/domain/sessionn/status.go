package sessionn

const (
	StatusNew      = "new"      // Новая
	StatusVerified = "verified" // Подтвержденная
	StatusExpired  = "expired"  // Истекшая
	StatusRevoked  = "revoked"  // Отозванная
)

func Statuses() []string {
	return []string{
		StatusNew,
		StatusVerified,
		StatusExpired,
		StatusRevoked,
	}
}
