package http2

import (
	"encoding/json"
	"fmt"
	"slices"
)

func DecodeBody(context Context, dst any) error {
	if err := json.NewDecoder(context.Request().Body).Decode(dst); err != nil {
		return fmt.Errorf("json Decoder Decode: %w", err)
	}

	return nil
}

func PathStr(context Context, name string) string {
	return context.Request().PathValue(name)
}

func QueryStr(context Context, name string) string {
	return context.Request().URL.Query().Get(name)
}

func Chain(h HandlerFuncRW, middlewares ...Middleware) HandlerFuncRW {
	for _, mw := range slices.Backward(middlewares) {
		h = mw(h)
	}
	return h
}
