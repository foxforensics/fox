package stream

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/cuhsat/fox/v4/internal"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

type Streamable interface {
	Write(*event.Event) error
}

type Stream struct {
	Tx  int64             `json:"-"`
	Rx  int64             `json:"-"`
	Url string            `json:"-"`
	Map map[string]string `json:"-"`
}

func (st *Stream) Post(body string) error {
	req, err := http.NewRequest("POST", st.Url, strings.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Add("user-agent", fmt.Sprintf("fox %s", app.Version))

	for k, v := range st.Map {
		req.Header.Set(k, v)
	}

	res, err := new(http.Client).Do(req)

	if err != nil {
		return err
	}

	st.Tx += req.ContentLength
	st.Rx += res.ContentLength

	if res.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(res.StatusCode))
	}

	return res.Body.Close()
}
