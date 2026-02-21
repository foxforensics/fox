package client

import (
	"fmt"
	"net/http"
	"time"

	res "github.com/cuhsat/fox/v4/internal"
)

var (
	Idle      = 0
	Timeout   = time.Second * 30
	UserAgent = fmt.Sprintf("fox %s", res.Version)
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
