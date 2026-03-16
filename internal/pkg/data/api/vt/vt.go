package vt

import (
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/VirusTotal/vt-go"

	"github.com/cuhsat/fox/v4/internal/pkg/data/api"
	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
)

var badCategories = []string{
	"malicious",
	"suspicious",
}

func CheckIp(ip, key string) (*api.Result, error) {
	return request(vt.URL("ip_addresses/%s", ip), key)
}

func CheckUrl(url, key string) (*api.Result, error) {
	return request(vt.URL("urls/%s", url), key)
}

func CheckDomain(url, key string) (*api.Result, error) {
	return request(vt.URL("domains/%s", url), key)
}

func CheckFileHash(sum, key string) (*api.Result, error) {
	return request(vt.URL("files/%s", sum), key)
}

func parseVerdict(obj *vt.Object, res *api.Result) {
	res.Stats.Bad = countStats(obj, badCategories)
	res.Stats.All = countStats(obj, []string{
		"malicious",
		"suspicious",
		"undetected",
		"harmless",
		"timeout",
		"confirmed-timeout",
		"failure",
		"type-unsupported",
	})

	res.Verdict, _ = obj.GetString("popular_threat_classification.suggested_threat_label")

	if len(res.Verdict) == 0 {
		switch {
		case res.Stats.Bad > 0:
			res.Verdict = api.Suspicious
		case res.Stats.All > 0:
			res.Verdict = api.Clean
		default:
			res.Verdict = api.Unrated
		}
	}
}

func parseDetails(obj *vt.Object, res *api.Result) {
	lar, err := obj.Get("last_analysis_results")

	if err != nil {
		return
	}

	m := lar.(map[string]any)

	for _, k := range slices.Sorted(maps.Keys(m)) {
		v := m[k].(map[string]any)

		if v["result"] == nil {
			continue
		}

		res.Details[v["engine_name"].(string)] = v["result"].(string)
	}
}

func countStats(obj *vt.Object, lst []string) (n int) {
	for _, k := range lst {
		v, _ := obj.GetInt64(fmt.Sprintf("last_analysis_stats.%s", k))
		n += int(v)
	}
	return
}

func request(url *url.URL, key string) (*api.Result, error) {
	res := &api.Result{Details: make(map[string]string)}

	vtc := vt.NewClient(key, vt.WithHTTPClient(client.Default()))

	obj, err := vtc.GetObject(url)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			res.Verdict = api.Unknown
			return res, nil
		}

		return nil, err
	}

	parseDetails(obj, res)
	parseVerdict(obj, res)

	return res, nil
}
