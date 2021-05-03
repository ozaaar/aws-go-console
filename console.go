// Package console provides short-lived (scoped based) token/url for AWS console
package console

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// Console provides the API operation methods for getting sign-in Token
type Console struct {
	STS    stsiface.STSAPI
	Client HTTPClient
}

// New creates a new instance of the Console Client with a session.
//
// Example:
//     mySession := session.Must(session.NewSession())
//
//     // Create a Console Client from just a session.
//     svc := sts.New(mySession)
func New(sess *session.Session) *Console {
	return &Console{
		STS:    sts.New(sess),
		Client: http.DefaultClient,
	}
}

// ErrInvalidToken indicates the token value is empty or is expired
var ErrInvalidToken = errors.New("invalid token")

// HTTPClient minimal interface for an HTTP client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Token contains the sign-in token for AWS console access
type Token struct {
	Value     string    `json:"SigninToken"`
	ExpiresAt time.Time `json:"-"`
}

// credentials used as payload for getting sign-in token
type credentials struct {
	SessionID    string `json:"sessionId"`
	SessionKey   string `json:"sessionKey"`
	SessionToken string `json:"sessionToken"`
}

// IsValid validates a given token
func (t *Token) IsValid() bool {
	if t.Value != "" && !t.ExpiresAt.IsZero() && t.ExpiresAt.After(time.Now()) {
		return true
	}
	return false
}

// SignInURL returns the URL with a valid token, can be opened directly in the browser
// dst is the AWS console location to be opened e.g https://console.aws.amazon.com/sns
func (t *Token) SignInURL(dst string) (*url.URL, error) {
	if !t.IsValid() {
		return nil, ErrInvalidToken
	}

	if _, err := url.ParseRequestURI(dst); err != nil {
		return nil, fmt.Errorf("invalid destination: %w", err)
	}

	rawUrl := fmt.Sprintf(
		"https://signin.aws.amazon.com/federation?Action=login&Destination=%s&SigninToken=%s",
		url.QueryEscape(dst),
		t.Value,
	)

	signInUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid sign-in url: %w", err)
	}

	return signInUrl, nil
}

// SignInTokenWithArn gets token from AWS API via GetFederationToken
// name is the name for AWS console user
// arn is the managed permission ARN for AWS console user
func (c *Console) SignInTokenWithArn(name, arn *string) (*Token, error) {
	input := sts.GetFederationTokenInput{
		Name:       name,
		PolicyArns: []*sts.PolicyDescriptorType{{Arn: arn}},
	}

	return c.signInToken(input)
}

// signInToken returns token against given credentials
func (c *Console) signInToken(input sts.GetFederationTokenInput) (*Token, error) {
	output, err := c.STS.GetFederationToken(&input)
	if err != nil {
		return nil, fmt.Errorf("getting federation token: %w", err)
	}

	cred := credentials{
		*output.Credentials.AccessKeyId,
		*output.Credentials.SecretAccessKey,
		*output.Credentials.SessionToken,
	}

	data, err := json.Marshal(cred)
	if err != nil {
		return nil, fmt.Errorf("marshalling session credentials: %w", err)
	}

	tokenRequestEndpoint := fmt.Sprintf(
		"https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=%s",
		url.QueryEscape(string(data)),
	)

	tokenRequest, err := http.NewRequest(http.MethodGet, tokenRequestEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating http request: %w", err)
	}

	tokenResponse, err := c.Client.Do(tokenRequest)
	defer tokenResponse.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("getting token response: %w", err)
	}

	tokenResponseBody, err := ioutil.ReadAll(tokenResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("reading token resposne: %w", err)
	}

	token := Token{ExpiresAt: time.Now()}
	err = json.Unmarshal(tokenResponseBody, &token)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling token: %w", err)
	}

	return &token, nil
}
