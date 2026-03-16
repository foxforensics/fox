package hibp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cuhsat/fox/v4/internal/pkg/data/api"
	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
)

const api3 = "https://haveibeenpwned.com/api/v3"

type breach struct {
	Name string `json:"Name,omitempty"`
}

func CheckMail(mail, key string) (*api.Result, error) {
	return request(fmt.Sprintf("%s/breachedaccount/%s?truncateResponse=false", api3, url.QueryEscape(mail)), key)
}

func parseVerdict(br []breach, res *api.Result) {
	res.Stats.All = len(br)
	res.Stats.Bad = len(br)

	if len(br) == 0 {
		res.Verdict = api.Clean
	} else {
		res.Verdict = api.Compromised
	}
}

func parseDetails(br []breach, res *api.Result) {
	for _, v := range br {
		res.Details[v.Name] = api.Compromised
	}
}

func getBreaches(resp *http.Response) ([]breach, error) {
	var br []breach

	b, err := io.ReadAll(resp.Body)

	_ = resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return br, json.Unmarshal(b, &br)
}

func request(url, key string) (*api.Result, error) {
	res := &api.Result{Details: make(map[string]string)}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("user-agent", client.UserAgent)
	req.Header.Add("hibp-api-key", key)

	resp, err := client.Default().Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	br, err := getBreaches(resp)

	if err != nil {
		return nil, err
	}

	parseDetails(br, res)
	parseVerdict(br, res)

	return res, nil
}
