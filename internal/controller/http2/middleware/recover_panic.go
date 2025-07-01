package middleware

import (
	"fmt"

	"github.com/nice-pea/npchat/internal/controller/http2"
)

// RecoverPanic перехватывает панику
func RecoverPanic(next http2.HandlerFuncRW) http2.HandlerFuncRW {
	return func(context http2.RWContext) (resp any, err error) {
		defer func() {
			if pv := recover(); pv != nil {
				resp = nil
				err = fmt.Errorf("%v", pv)
				return
			}
		}()

		return next(context)
	}
}
