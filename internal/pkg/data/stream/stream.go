package stream

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

var userAgent = fmt.Sprintf("fox %s", app.Version)
var httpClient = new(http.Client)

type Streamable interface {
	Write(*event.Event) error
}

type Streamer struct {
	Map map[string]string `json:"-"`
	Url string            `json:"-"`
}

func (str *Streamer) Post(body string) error {
	req, err := http.NewRequest("POST", str.Url, strings.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("user-agent", userAgent)

	for k, v := range str.Map {
		req.Header.Set(k, v)
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(res.StatusCode))
	}

	return res.Body.Close()
}
