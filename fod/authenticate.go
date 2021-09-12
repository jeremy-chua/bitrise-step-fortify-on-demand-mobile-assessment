package fod

import (
	"errors"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// Client Credentials ...
type ClientCredentials struct {
	ClientID     string `structs:"client_id"`
	ClientSecret string `structs:"client_secret"`
	GrantType    string `structs:"grant_type"`
	Scope        string `structs:"scope"`
}

// User Credentials ...
type UserCredentials struct {
	Username  string `structs:"username"`
	Password  string `structs:"password"`
	GrantType string `structs:"grant_type"`
	Scope     string `structs:"scope"`
}

// Auth Data ...
type AuthData struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// Auth Error ...
type AuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (c *Client) authenticate() (*AuthData, error) {

	var (
		resp       *resty.Response = nil
		err        error           = nil
		authData                   = &AuthData{}
		url                        = c.baseUrl + post_api_authenticate
		authParams interface{}
	)

	switch c.grantType {
	case grant_type_client_credentials:
		authParams = ClientCredentials{
			ClientID:     c.username,
			ClientSecret: c.password,
			GrantType:    c.grantType,
			Scope:        scopeToString(c.scope),
		}
	case grant_type_password:
		authParams = UserCredentials{
			Username:  c.username,
			Password:  c.password,
			GrantType: c.grantType,
			Scope:     scopeToString(c.scope),
		}
	default:
		log.Error("invalid invalid grant type")
		return nil, errors.New("invalid grant type")
	}

	log.WithFields(log.Fields{
		"grant_type": c.grantType,
		"base_url":   c.baseUrl,
	}).Info("authenticating...")

	if resp, err = c.webClient.R().
		SetFormData(structToFormData(authParams)).
		SetResult(authData).
		SetError(&AuthError{}).
		Post(url); err == nil {

		switch resp.StatusCode() {
		case 200:
			log.Info("authenticated")
			return authData, nil
		default:
			log.WithFields(log.Fields{
				"status_code":       resp.StatusCode(),
				"error":             resp.Error().(*AuthError).Error,
				"error_description": resp.Error().(*AuthError).ErrorDescription,
			}).Error("unsuccessful authentication")
			err = errors.New(resp.Error().(*AuthError).ErrorDescription)
		}
	}

	return nil, err
}

func (c *Client) TestAuthenticate(scope ...string) (string, error) {

	c.scope = scope
	authData, err := c.authenticate()
	return authData.AccessToken, err
}
