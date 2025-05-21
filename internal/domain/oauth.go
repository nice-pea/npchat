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
}

type OAuthGoogleUser struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
}

type OAuthLink struct {
	State      string
	UserID     string
	ExternalID string
	//	Provider string
}

type OAuthRepository interface {
	SaveToken(OAuthToken) error
	SaveLink(OAuthLink) error
	Link(OAuthLinkFilter) ([]OAuthLink, error)
}

type OAuthLinkFilter struct {
	State      string
	UserID     string
	ExternalID string
	//	Provider string
}
