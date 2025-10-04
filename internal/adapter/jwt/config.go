package jwt

type Config struct {
	SecretKey                   string // Ключ для подписи JWT
	VerifyTokenWithInvalidation bool   // использовать ли проверку анулирования токена
	RedisDSN                    string // redis DSN
}
