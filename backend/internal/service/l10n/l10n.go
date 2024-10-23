package l10n

type Service interface {
	Localize(code, locale string, vars map[string]string) (string, error)
}
