package http

import (
	"errors"
	"net/http"
	"strings"

	"go.foxforensics.dev/fox/v4/internal/pkg/types/client"
)

var _client = client.Http()

func Post(url, body string, headers map[string]string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("user-agent", client.Name)

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
