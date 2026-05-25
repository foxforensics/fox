package client

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.golang/paho"

	"go.foxforensics.dev/fox/v4/internal/pkg/version"
)

var (
	Name    = fmt.Sprintf("fox %s", version.Number)
	Timeout = time.Second * 30
	MaxIdle = 0
)

// Http return the default http client
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

// Mqtt returns the default mqtt client
func Mqtt(url, usr, pwd string) *mqtt.Client {
	conn, err := net.Dial("tcp", strings.TrimPrefix(url, "tcp://"))

	if err != nil {
		log.Fatalln(err)
	}

	c := mqtt.NewClient(mqtt.ClientConfig{Conn: conn})

	cp := &mqtt.Connect{
		KeepAlive:  uint16(Timeout.Seconds()),
		ClientID:   Name,
		CleanStart: true,
	}

	if len(usr) > 0 {
		cp.Username = usr
		cp.UsernameFlag = true
	}

	if len(pwd) > 0 {
		cp.Password = []byte(pwd)
		cp.PasswordFlag = true
	}

	ca, err := c.Connect(context.Background(), cp)

	if err != nil {
		log.Fatalln(err)
	}

	if ca.ReasonCode > 0 {
		log.Fatalln(ca.Properties.ReasonString)
	}

	return c
}
