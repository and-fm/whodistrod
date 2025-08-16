package dspclients

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	apiclient "github.com/and-fm/whodistrod/internal/api_client"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type tidalClient struct {
	deps TidalClientDeps

	encodedClientCreds string
	tokenUrl           string
	currentToken       TidalToken
}

type TidalClientDeps struct {
	dig.In

	Logger    logging.Logger
	ApiClient apiclient.ApiClient
}

type TidalClient interface {
	Authenticate() error
	ApiClient() apiclient.ApiClient
	DoRequest(req *resty.Request) (*resty.Response, error)
}

type TidalClientConfig struct {
	ClientId     string
	ClientSecret string
	TokenURL     string
}

type TidalToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewTidalClient(conf TidalClientConfig, d TidalClientDeps) TidalClient {
	encodedClientCreds := base64.StdEncoding.EncodeToString([]byte(conf.ClientId + ":" + conf.ClientSecret))

	client := tidalClient{
		deps:               d,
		encodedClientCreds: encodedClientCreds,
		tokenUrl:           conf.TokenURL,
	}

	return &client
}

func (c *tidalClient) Authenticate() error {
	res, err := c.AuthWithClientCredentials()

	if err != nil {
		return fmt.Errorf("auth with tidal client creds: %w", err)
	}

	token, err := GetTidalTokenObjectFromBody(res.Body())

	if err != nil {
		return fmt.Errorf("unmarshalling tidal token response: %w", err)
	}

	c.currentToken = token

	c.deps.Logger.LogInfo("authenticated with tidal using client credentials", nil)

	return nil
}

func (c *tidalClient) AuthWithClientCredentials() (*resty.Response, error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req := c.deps.ApiClient.NewRequest(c.tokenUrl)

	req.FormData = form

	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	req.SetHeader("Authorization", "Basic "+c.encodedClientCreds)

	res, err := req.Execute("POST", req.URL)

	return res, err
}

func (c *tidalClient) ApiClient() apiclient.ApiClient {
	return c.deps.ApiClient
}

func GetTidalTokenObjectFromBody(body []byte) (TidalToken, error) {
	token := TidalToken{}

	err := json.Unmarshal(body, &token)
	return token, err
}

func (c *tidalClient) DoRequest(req *resty.Request) (*resty.Response, error) {
	req.SetHeader("Authorization", "Bearer "+c.currentToken.AccessToken)
	req.SetHeader("accept", "application/vnd.api+json")

	res, err := req.Execute(req.Method, req.URL)

	if err != nil {
		return nil, fmt.Errorf("executing tidal request: %w", err)
	}

	if res.StatusCode() > 399 {
		return nil, echo.NewHTTPError(res.StatusCode(), "Tidal request failed")
	}

	return res, nil
}
