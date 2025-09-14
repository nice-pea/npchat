package register_handler

type JwtIssuer interface {
	Issue(claims map[string]any) (string, error)
}
