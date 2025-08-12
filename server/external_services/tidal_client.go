package externsvc

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/labstack/echo/v4"
)

type tidalClient struct {
	client             http.Client
	logger logging.Logger
	encodedClientCreds string
	tokenUrl           string
	currentToken       TidalToken
}

type TidalClient interface {
	Authenticate() error
	DoRequest(req *http.Request) (*http.Response, error)
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

func NewTidalClient(conf TidalClientConfig, logger logging.Logger) TidalClient {
	encodedClientCreds := base64.StdEncoding.EncodeToString([]byte(conf.ClientId + ":" + conf.ClientSecret))

	client := tidalClient{
		client:             http.Client{},
		encodedClientCreds: encodedClientCreds,
		tokenUrl:           conf.TokenURL,
		logger: logger,
	}

	return &client
}

func (c *tidalClient) Authenticate() error {
	res, err := c.AuthWithClientCredentials()

	if err != nil {
		return err
	}

	token, err := GetTokenObjectFromBody(res.Body)

	if err != nil {
		return err
	}

	c.currentToken = token

	c.logger.LogInfo("authenticated with tidal using client credentials", nil)

	return nil
}

func (c *tidalClient) AuthWithClientCredentials() (*http.Response, error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", c.tokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+c.encodedClientCreds)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.client.Do(req)

	return res, err
}

func GetTokenObjectFromBody(body io.ReadCloser) (TidalToken, error) {
	defer body.Close()

	token := TidalToken{}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(bodyBytes, &token)
	return token, err
}

func (c *tidalClient) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.currentToken.AccessToken)
	req.Header.Add("accept", "application/vnd.api+json")
	res, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode > 399 {
		return nil, echo.NewHTTPError(res.StatusCode)
	}

	return res, nil
}

func Do[T any](tidalClient TidalClient, method string, endpoint string, body any) (T, error) {
	var obj T

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return obj, err
	}
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return obj, err
		}
		reader := bytes.NewReader(b)
		req, err = http.NewRequest(method, endpoint, reader)
		if err != nil {
			return obj, err
		}
	}

	res, err := tidalClient.DoRequest(req)

	if err != nil {
		return obj, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return obj, err
	}
	fmt.Println(string(resBody))

	err = json.Unmarshal(resBody, &obj)
	return obj, err
}
