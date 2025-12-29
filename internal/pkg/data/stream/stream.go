package stream

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

// statics
var (
	agent  = fmt.Sprintf("fox %s", app.Version)
	client = http.Client{
		Timeout: time.Second * 30,
	}
)

type Streamer interface {
	Stream(*event.Event) error
}

func Post(url, body string, headers map[string]string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("user-agent", agent)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(res.StatusCode))
	}

	return res.Body.Close()
}
