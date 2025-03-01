package controller

type Localization interface {
	Localize(code, locale string, vars map[string]string) (string, error)
}

type Config struct {
}
