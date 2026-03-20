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

// Encrypted reserve keys for emergency use
const (
	ReserveKey1 = "47ba3c085f105fff4fa186ce769f8a35f98bc3010fd8e25c9a90c1bf70696120b9fe1a5c6328bf0deae4eebdcc9f5df156a27efd923eaad648f3e8ab26fcc8f6753233b8"
	ReserveKey2 = "44201ef4cbffe7edd1a7d2279a1fc3019700c3620da45d0542014b8a7be0fd7b53125c3e474c6db7360f4f538d56bfe15bd416b0d2a77c02a37d0ffc5015694b41c9f117"
)

var badCategories = []string{
	"malicious",
	"suspicious",
}

func CheckFile(sum, key string) (*api.Result, error) {
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
