package domain

import "time"

type OAuthToken struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time

	// LinkID is the ID of the link associated with the token.
	// TODO: Required.
	LinkID string

	Provider string
}

type OAuthUser struct {
	ID       string
	Email    string
	Name     string
	Picture  string
	Provider string
}

type OAuthLink struct {
	ID         string
	UserID     string
	ExternalID string
	Provider   string
}

type OAuthRepository interface {
	SaveToken(OAuthToken) error
	SaveLink(OAuthLink) error
	ListLinks(OAuthListLinksFilter) ([]OAuthLink, error)
}

type OAuthListLinksFilter struct {
	ID         string
	UserID     string
	ExternalID string
	Provider   string
}
