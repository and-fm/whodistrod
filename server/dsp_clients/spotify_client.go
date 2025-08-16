package dspclients

import (
	"encoding/json"
	"fmt"
	"net/url"

	apiclient "github.com/and-fm/whodistrod/internal/api_client"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type spotifyClient struct {
	deps         SpotifyClientDeps
	clientConfig SpotifyClientConfig
	currentToken SpotifyToken
}

type SpotifyClientDeps struct {
	dig.In

	Logger    logging.Logger
	ApiClient apiclient.ApiClient
}

type SpotifyClient interface {
	Authenticate() error
	ApiClient() apiclient.ApiClient
	DoRequest(req *resty.Request) (*resty.Response, error)
}

type SpotifyClientConfig struct {
	ClientId     string
	ClientSecret string
	TokenURL     string
}

type SpotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewSpotifyClient(conf SpotifyClientConfig, d SpotifyClientDeps) SpotifyClient {
	client := spotifyClient{
		clientConfig: conf,
		deps:         d,
	}

	return &client
}

func (c *spotifyClient) Authenticate() error {
	res, err := c.AuthWithClientCredentials()

	if err != nil {
		return fmt.Errorf("auth with spotify client creds: %w", err)
	}

	token, err := GetSpotifyTokenObjectFromBody(res.Body())

	if err != nil {
		return fmt.Errorf("unmarshalling spotify token response: %w", err)
	}

	c.currentToken = token

	c.deps.Logger.LogInfo("authenticated with spotify using client credentials", nil)

	return nil
}

func (c *spotifyClient) AuthWithClientCredentials() (*resty.Response, error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", c.clientConfig.ClientId)
	form.Add("client_secret", c.clientConfig.ClientSecret)

	req := c.deps.ApiClient.NewRequest(c.clientConfig.TokenURL)

	req.FormData = form

	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	res, err := req.Execute("POST", req.URL)

	return res, err
}

func (c *spotifyClient) ApiClient() apiclient.ApiClient {
	return c.deps.ApiClient
}

func GetSpotifyTokenObjectFromBody(body []byte) (SpotifyToken, error) {
	token := SpotifyToken{}

	err := json.Unmarshal(body, &token)
	return token, err
}

func (c *spotifyClient) DoRequest(req *resty.Request) (*resty.Response, error) {
	req.SetHeader("Authorization", "Bearer "+c.currentToken.AccessToken)
	res, err := req.Execute(req.Method, req.URL)

	if err != nil {
		return nil, fmt.Errorf("executing spotify request: %w", err)
	}

	if res.StatusCode() > 399 {
		return nil, echo.NewHTTPError(res.StatusCode(), "Spotify request failed")
	}

	return res, nil
}
