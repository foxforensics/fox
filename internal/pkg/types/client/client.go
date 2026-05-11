package client

import (
	"fmt"
	"net/http"
	"time"

	"go.foxforensics.dev/fox/v4/internal"
)

var (
	Idle      = 0
	Timeout   = time.Second * 30
	UserAgent = fmt.Sprintf("fox %s", version.Number)
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
