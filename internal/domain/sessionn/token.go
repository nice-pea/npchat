package sessionn

import (
	"time"
)

type Token struct {
	Token  string
	Expiry time.Time
}
