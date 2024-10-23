package l10n

type Service interface {
	Localize(category, name, locale string, params map[string]string) (string, error)
}
