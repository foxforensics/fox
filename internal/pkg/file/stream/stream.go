package stream

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

var _client = client.Default()

type Streamer interface {
	Stream(*event.Event) error
}

func Post(url, body string, headers map[string]string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("user-agent", client.UserAgent)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := _client.Do(req)

	defer func() {
		_ = res.Body.Close()
	}()

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(res.StatusCode))
	}

	return nil
}
