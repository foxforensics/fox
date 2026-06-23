package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.golang/paho"
	"go.foxforensics.eu/fox/v4/internal/sys/version"

	"github.com/segmentio/ksuid"
)

// Timeout for everything network related.
var Timeout = time.Second * 30

// MaxIdle connections at once.
var MaxIdle = 0

func ID() string {
	uid, err := ksuid.NewRandomWithTime(time.Now().UTC())

	if err != nil {
		slog.Error(err.Error())
		return Name()
	}

	return fmt.Sprintf("%s %s", Name(), uid.String())
}

// Name returns the clients name including the version number.
func Name() string {
	return fmt.Sprintf("fox %s", version.Number)
}

// Http return the default http client.
func Http() *http.Client {
	return &http.Client{
		Timeout: Timeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			IdleConnTimeout:     Timeout,
			TLSHandshakeTimeout: Timeout,
			MaxIdleConnsPerHost: MaxIdle,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12, // pinned
			},
		},
	}
}

// Mqtt returns the default mqtt client.
func Mqtt(adr, usr, pwd string) (*mqtt.Client, error) {
	u, err := url.Parse(adr)

	if err != nil {
		return nil, err
	}

	con, err := tls.Dial("tcp", u.Host, &tls.Config{
		MinVersion: tls.VersionTLS12, // pinned
	})

	if err != nil {
		return nil, err
	}

	client := mqtt.NewClient(mqtt.ClientConfig{Conn: con})

	pkg := &mqtt.Connect{
		KeepAlive:  uint16(Timeout.Seconds()),
		ClientID:   ID(),
		CleanStart: true,
	}

	if len(usr) > 0 {
		pkg.Username = usr
		pkg.UsernameFlag = true
	}

	if len(pwd) > 0 {
		pkg.Password = []byte(pwd)
		pkg.PasswordFlag = true
	}

	ctx, stop := context.WithTimeout(context.Background(), Timeout)
	defer stop()

	ack, err := client.Connect(ctx, pkg)

	if err != nil {
		return nil, err
	}

	if ack.ReasonCode > 0 {
		return nil, errors.New(ack.Properties.ReasonString)
	}

	return client, nil
}
