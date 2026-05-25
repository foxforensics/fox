package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.golang/paho"

	"go.foxforensics.dev/fox/v4/internal/pkg/version"
)

// Timeout for everything network related.
var Timeout = time.Second * 30

// MaxIdle connections at once.
var MaxIdle = 0

// ID returns the clients unique and reproducible id. It is build from
// the programs name and version number, followed by the SHA256 hash of
// the hostname and the first interface found that is up.
//
// Example:
//
//	fox 1.2.3 42fee663a1683b00383ec69d91e4880335cd6b265611e4e7b4cdf5e5e4ae22d7
func ID() string {
	id, err := os.Hostname()

	if err != nil {
		log.Fatalln(err)
	}

	iff, err := net.Interfaces()

	if err != nil {
		log.Fatalln(err)
	}

	for _, i := range iff {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			id = fmt.Sprintf("%s %s", id, i.HardwareAddr.String())
			break
		}
	}

	return fmt.Sprintf("%s %x", Name(), sha256.Sum256([]byte(id)))
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
		},
	}
}

// Mqtt returns the default mqtt client.
func Mqtt(url, usr, pwd string) *mqtt.Client {
	conn, err := net.Dial("tcp", strings.TrimPrefix(url, "tcp://"))

	if err != nil {
		log.Fatalln(err)
	}

	c := mqtt.NewClient(mqtt.ClientConfig{Conn: conn})

	cp := &mqtt.Connect{
		KeepAlive:  uint16(Timeout.Seconds()),
		ClientID:   ID(),
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
