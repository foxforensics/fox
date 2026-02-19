package client

import (
	"fmt"
	"net/http"
	"time"

	app "github.com/cuhsat/fox/v4/internal"
)

var (
	Idle      = 0
	Timeout   = time.Second * 30
	UserAgent = fmt.Sprintf("fox %s", app.Version)
)

// Default HTTP client
func Default() *http.Client {
	return &http.Client{
		Timeout: Timeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			IdleConnTimeout:     Timeout,
			TLSHandshakeTimeout: Timeout,
			MaxIdleConnsPerHost: Idle,
		},
	}
}
