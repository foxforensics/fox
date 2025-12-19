package stream

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

var userAgent = fmt.Sprintf("fox %s", app.Version)
var httpClient = new(http.Client)

type Streamable interface {
	Write(*event.Event) (int64, int64, error)
}

type Stream struct {
	Url string            `json:"-"`
	Map map[string]string `json:"-"`
}

func (stm *Stream) Post(body string) (int64, int64, error) {
	req, err := http.NewRequest("POST", stm.Url, strings.NewReader(body))

	if err != nil {
		return 0, 0, err
	}

	req.Header.Add("user-agent", userAgent)

	for k, v := range stm.Map {
		req.Header.Set(k, v)
	}

	res, err := httpClient.Do(req)

	if err != nil {
		return 0, 0, err
	}

	tx := req.ContentLength

	buf, err := io.ReadAll(res.Body)

	if err != nil {
		return tx, 0, err
	}

	rx := int64(len(buf))

	if res.StatusCode != http.StatusOK {
		return tx, rx, errors.New(http.StatusText(res.StatusCode))
	}

	return tx, rx, res.Body.Close()
}
