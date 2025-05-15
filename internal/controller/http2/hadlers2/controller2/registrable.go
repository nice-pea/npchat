package controller2

type Registrable interface {
	RegisterHandlers(...Handler2) error
}
