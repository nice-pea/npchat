package jwt

type Config struct {
	SecretKey                     string // Ключ для подписи JWT
	VerifyTokenWithAdvancedChecks bool   // использовать ли продвинутую проверку токена
	RedisDSN                      string // redis DSN
}
