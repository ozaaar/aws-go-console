package console

import (
	"net/url"
	"time"
)

// Token contains the sign-in token for AWS console access
type Token struct {
	Value     string    `json:"SigninToken"`
	ExpiresAt time.Time `json:"-"`
}

// SignInToken fetches the token from AWS API via GetFederationToken
func SignInToken() (*Token, error) {
	panic("not implemented")
}

// IsValid confirms the validity for a given token
func (t *Token) IsValid() bool {
	if t.Value != "" && !t.ExpiresAt.IsZero() && t.ExpiresAt.After(time.Now()) {
		return true
	}
	return false
}

// SignInURL returns the URL with token included for short time access
func (t *Token) SignInURL() *url.URL {
	panic("not implemented")
}
