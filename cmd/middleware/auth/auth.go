package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// UserInfo contains a user's basic information extracted from an access token.
type UserInfo struct {
	SubID     string `json:"sub,omitempty"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Picture   string `json:"picture"`
	Email     string `json:"email"`
}

// Config contains the information required to configure a validator to make
// requests to validate a request's authorization header.
type Config struct {
	// Domain represents the domain where the user info endpoint is hosted.
	Domain string
	// BasicAuthCredentials represents the base64 encoded username and password
	// used to authenticate another service with basic auth.
	BasicAuthCredentials string
}

// Validate looks at the configuration's contents to ensure it has all the
// required fields.
func (conf *Config) validate() error {
	if conf.Domain == "" {
		return errors.New("missing domain")
	}

	return nil
}

// Validator is an interface representing the ability to validate an
// authorization header and return the authenticated user's information.
type Validator interface {
	// Validate validates an authorization and returns the authenticated user's
	// information.
	Validate(credentials string) (*UserInfo, error)
}

// A TokenValidator is a validator that validates a bearer token in an
// authorization header by making a request to a /userinfo endpoint.
type TokenValidator struct {
	conf *Config
}

// NewTokenValidator creates a new token validator with the given
// configuration.
func NewTokenValidator(conf *Config) (Validator, error) {
	if conf == nil {
		return nil, fmt.Errorf("auth: missing configuration")
	}

	err := conf.validate()
	if err != nil {
		return nil, fmt.Errorf("auth: configuration %s", err)
	}

	return &TokenValidator{conf}, nil
}

// Validate makes a request to the /userinfo endpoint on the domain specified
// in the token validator's configuration to validate the bearer token present
// in the authorization header and returns the authenticated user's
// information.
func (validator *TokenValidator) Validate(credentials string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://"+validator.conf.Domain+"/userinfo", nil)
	if err != nil {
		return nil, UnauthorizedError{fmt.Sprintf("auth.TokenValidator: failed to create request (%s)", err)}
	}

	req.Header.Set("Authorization", "Bearer "+credentials)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, UnauthorizedError{fmt.Sprintf("auth: failed to make request (%s)", err)}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, UnauthorizedError{fmt.Sprintf("auth: failed to validate token")}
	}

	var userInfo UserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, UnauthorizedError{fmt.Sprintf("auth: failed to decode user info (%s)", err)}
	}
	return &userInfo, nil
}

// A BasicAuthValidator is a validator that validates a base64 encoded username
// and password is the same as the one that is configured.
type BasicAuthValidator struct {
	conf *Config
}

// NewBasicAuthValidator creates a new token validator with the given
// configuration.
func NewBasicAuthValidator(conf *Config) (Validator, error) {
	if conf == nil {
		return nil, fmt.Errorf("auth: missing configuration")
	}

	err := conf.validate()
	if err != nil {
		return nil, fmt.Errorf("auth: configuration %s", err)
	}

	return &BasicAuthValidator{conf}, nil
}

// Validate compares the authorization header with the base64 encoded username
// and password stored in its configuration. It does not return the
// authenticated user's information, since there is no user.
func (validator *BasicAuthValidator) Validate(credentials string) (*UserInfo, error) {
	if strings.Compare(credentials, validator.conf.BasicAuthCredentials) == 0 {
		return &UserInfo{}, nil
	}

	return nil, UnauthorizedError{"auth: failed to decode user info"}
}

type contextKey string

func (c contextKey) String() string {
	return "auth." + string(c)
}

const (
	// UserInfoContextKey represents the key used to store and retrieve the
	// user information from the request context.
	UserInfoContextKey = contextKey("userInfo")
)

// FromContext extracts an authenticated user's information from the request's
// context.
func FromContext(ctx context.Context) (*UserInfo, error) {
	if ctx == nil {
		return nil, fmt.Errorf("auth: request context is nil")
	}

	userInfo, ok := ctx.Value(UserInfoContextKey).(*UserInfo)
	if !ok {
		return nil, fmt.Errorf("auth: %s not found in context", UserInfoContextKey)
	}

	return userInfo, nil
}
