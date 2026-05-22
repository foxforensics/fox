package client

import (
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.foxforensics.dev/fox/v4/internal/pkg/version"
)

var (
	Name    = fmt.Sprintf("fox %s", version.Number)
	Timeout = time.Second * 30
	MaxIdle = 0
)

// Http client
func Http() *http.Client {
	return &http.Client{
		Timeout: Timeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			IdleConnTimeout:     Timeout,
			TLSHandshakeTimeout: Timeout,
			MaxIdleConnsPerHost: MaxIdle,
		},
	}
}

// Mqtt client
func Mqtt(url string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID(Name)
	opts.SetPingTimeout(Timeout)
	opts.SetWriteTimeout(Timeout)
	return mqtt.NewClient(opts)
}
