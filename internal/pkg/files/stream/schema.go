package stream

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cuhsat/fox/v4/internal"
)

type Schema struct {
	Url string            `json:"-"`
	Map map[string]string `json:"-"`
}

func (sch *Schema) Post(s string) (int, error) {
	req, err := http.NewRequest("POST", sch.Url, strings.NewReader(s))

	if err != nil {
		return 0, err
	}

	req.Header.Add("user-agent", fmt.Sprintf("fox %s", app.Version))

	for k, v := range sch.Map {
		req.Header.Set(k, v)
	}

	res, err := new(http.Client).Do(req)

	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		return 0, errors.New(http.StatusText(res.StatusCode))
	}

	err = res.Body.Close()

	if err != nil {
		return 0, err
	}

	return len(s), nil
}
