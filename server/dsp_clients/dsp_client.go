package dspclients

import (
	apiclient "github.com/and-fm/whodistrod/internal/api_client"
	"github.com/go-resty/resty/v2"
)

type DspClient interface {
	Authenticate() error
	ApiClient() apiclient.ApiClient
	DoRequest(req *resty.Request) (*resty.Response, error)
}