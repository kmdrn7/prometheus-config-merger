package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

var httpClient *resty.Client

// Initialize http rest client
func New() {

	isDebug := viper.GetBool("debug")

	client := resty.New()
	client.SetDebug(isDebug)
	client.SetHeader("User-Agent", "prometheus-config-merger/1.0.0")
	httpClient = client
}

// Send http post request to specified url
func Post(url string) (*resty.Response, error) {
	r, err := httpClient.R().Post(url)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Send http get request to specified url
func Get(url string) {
	httpClient.R().Get(url)
}
