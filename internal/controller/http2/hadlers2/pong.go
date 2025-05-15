package hadlers2

import (
	"github.com/saime-0/nice-pea-chat/internal/controller"
)

type Pong struct{}

func (p *Pong) HandlerParams() controller.HandlerParams {
	return controller.HandlerParams{
		Method: "",
		Path:   "/ping",
	}
}

func (p *Pong) HandlerFunc(context controller.Context) (any, error) {
	return "pong", nil
}
