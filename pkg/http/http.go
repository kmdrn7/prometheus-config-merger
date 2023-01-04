package http

import "github.com/go-resty/resty/v2"

var httpClient *resty.Client

// Initialize http rest client
func New() {
	client := resty.New()
	client.SetDebug(true)
	httpClient = client
}

// Send http post request to specified url
func Post(url string) {
	httpClient.R().Post(url)
}

// Send http get request to specified url
func Get(url string) {
	httpClient.R().Get(url)
}
