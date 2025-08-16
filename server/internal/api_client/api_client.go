package apiclient

import (
	"encoding/json"

	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type apiClient struct {
	dig.In `ignore-unexported:"true"`

	client *resty.Client // not injected, it is created in NewApiClient

	Logger logging.Logger
}

type ApiClient interface {
	NewRequestWithBody(url string, body any) *resty.Request
	NewRequest(url string) *resty.Request
}

func NewApiClient(a apiClient) ApiClient {
	client := resty.New()

	client.SetRetryCount(3)

	restyLogger := NewApiClientLogger(a.Logger)

	client.SetLogger(restyLogger)

	client.OnError(func(r *resty.Request, err error) {
		if v, ok := err.(*resty.ResponseError); ok {
			a.Logger.LogError("Request failed with status code %d: %s - %s", nil, v.Response.StatusCode(), v.Error(), err.Error())
		} else {
			a.Logger.LogError("Unknown request error: %s ", nil, err.Error())
		}
	})

	a.client = client

	return &a
}

func (a *apiClient) Post(host string, endpoint string, body any, result any) (*resty.Response, error) {
	res, err := a.client.R().SetBody(body).Post(host + endpoint)

	return res, err
}

func (a *apiClient) NewRequestWithBody(url string, body any) *resty.Request {
	req := a.client.R().SetBody(body)

	req.URL = url

	return req
}

func (a *apiClient) NewRequest(url string) *resty.Request {
	req := a.client.R()

	req.URL = url

	return req
}

func GetResponseError(res *resty.Response, err error) error {
	if res.StatusCode() > 399 {
		return echo.NewHTTPError(res.StatusCode())
	}

	return err
}

func GetResponseBody[T any](body []byte) (T, error) {
	var obj T

	err := json.Unmarshal(body, &obj)

	return obj, err
}

func GET[T any](req *resty.Request) (*resty.Response, T, error) {
	var obj T
	res, err := req.Get(req.URL)

	err = GetResponseError(res, err)
	if err != nil {
		return nil, obj, err
	}

	resBody, err := GetResponseBody[T](res.Body())

	return res, resBody, err
}

func POST[T any](req *resty.Request) (*resty.Response, T, error) {
	var obj T

	res, err := req.Post(req.URL)

	err = GetResponseError(res, err)
	if err != nil {
		return nil, obj, err
	}

	resBody, err := GetResponseBody[T](res.Body())

	return res, resBody, err
}

func PUT[T any](req *resty.Request) (*resty.Response, T, error) {
	var obj T

	res, err := req.Put(req.URL)

	err = GetResponseError(res, err)
	if err != nil {
		return nil, obj, err
	}

	resBody, err := GetResponseBody[T](res.Body())

	return res, resBody, err
}

func DELETE[T any](req *resty.Request) (*resty.Response, T, error) {
	var obj T

	res, err := req.Get(req.URL)

	err = GetResponseError(res, err)
	if err != nil {
		return nil, obj, err
	}

	resBody, err := GetResponseBody[T](res.Body())

	return res, resBody, err
}

