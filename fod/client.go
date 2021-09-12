package fod

import (
	"strings"

	"github.com/go-resty/resty/v2"

	log "github.com/sirupsen/logrus"
)

// Fortify on Demand RESTFul Client struct ...
type Client struct {
	username  string
	password  string
	baseUrl   string
	grantType string
	scope     []string
	webClient *resty.Client
}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

// Set debug to print debug information
func (c *Client) SetDebug(enable bool) *Client {

	if enable {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return c
}

// Creates new Fortify on Demand RESTFul client wih Client Credentials
func NewWithClientCredentials(clientId, secret, datacenter string) *Client {

	initLogger()

	// return nil for invalid arguments
	if strings.TrimSpace(clientId) == "" || strings.TrimSpace(secret) == "" || !isValidDatacenter(datacenter) {
		log.Error("invalid parameters when creating with client credentials")
		return nil
	}

	c := &Client{
		username:  clientId,
		password:  secret,
		baseUrl:   getBaseUrl(datacenter),
		grantType: grant_type_client_credentials,
		scope:     nil,
		webClient: resty.New().SetHeader("Accept", "application/json"),
	}

	return c
}

// Creates new Fortify on Demand RESTFul client with User Credentials
func NewWithUserCredentials(tenant, username, password, datacenter string) *Client {

	initLogger()

	// return nil for invalid arguments
	if strings.TrimSpace(tenant) == "" || strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" || !isValidDatacenter(datacenter) {
		log.Error("invalid parameters when creating with user credentials")
		return nil
	}

	c := &Client{
		username:  tenant + "\\" + username,
		password:  password,
		baseUrl:   getBaseUrl(datacenter),
		grantType: grant_type_password,
		scope:     nil,
		webClient: resty.New().SetHeader("Accept", "application/json"),
	}

	return c
}
