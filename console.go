package console

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// InvalidTokenError indicates the token value is empty or is expired
var InvalidTokenError = errors.New("invalid token")

// Token contains the sign-in token for AWS console access
type Token struct {
	Value     string    `json:"SigninToken"`
	ExpiresAt time.Time `json:"-"`
}

// SignInToken fetches the token from AWS API via GetFederationToken
func SignInToken() (*Token, error) {
	panic("not implemented")
}

// IsValid validates  a given token
func (t *Token) IsValid() bool {
	if t.Value != "" && !t.ExpiresAt.IsZero() && t.ExpiresAt.After(time.Now()) {
		return true
	}
	return false
}

// SignInURL returns the URL with a valid token
func (t *Token) SignInURL(dst string) (*url.URL, error) {
	if !t.IsValid() {
		return nil, InvalidTokenError
	}

	if _, err := url.ParseRequestURI(dst); err != nil {
		return nil, fmt.Errorf("invalid destination: %w", err)
	}

	rawurl := fmt.Sprintf(
		"https://signin.aws.amazon.com/federation?Action=login&Destination=%s&SigninToken=%s",
		url.QueryEscape(dst),
		t.Value,
	)

	signInUrl, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, fmt.Errorf("invalid sign-in url: %w", err)
	}

	return signInUrl, nil
}
