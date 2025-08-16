package dspclients

import (
	"errors"
	"fmt"

	apiclient "github.com/and-fm/whodistrod/internal/api_client"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

func Do[T any](dspClient DspClient, method string, endpoint string, body any) (*resty.Response, T, error) {
	var resBody T

	apiClient := dspClient.ApiClient()

	req := apiClient.NewRequest(endpoint)

	if body != nil {
		req.SetBody(body)
	}

	req.Method = method

	res, err := dspClient.DoRequest(req)

	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		if httpError.Code == 401 {
			err = dspClient.Authenticate()
			if err != nil {
				return nil, resBody, err
			}
			res, err = dspClient.DoRequest(req)
		} else {
			return nil, resBody, fmt.Errorf("dspclient request failed: %w", err)
		}
	}

	if err != nil {
		return nil, resBody, err
	}

	resBody, err = apiclient.GetResponseBody[T](res.Body())

	return res, resBody, err
}
