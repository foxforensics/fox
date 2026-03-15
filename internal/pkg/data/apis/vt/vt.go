package vt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/VirusTotal/vt-go"

	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
)

const (
	Unknown    = "unknown"
	Unrated    = "unrated"
	Clean      = "clean"
	Indecisive = "indecisive"
)

// Verbose API responses
var Verbose int

var alerts = []string{
	"malicious",
	"suspicious",
}

type Entry struct {
	Alert  bool
	Engine string
	Result string
}

type Result struct {
	Entries []Entry
	Label   string
	Alert   bool
	All     int
	Bad     int
}

func CheckIp(ip, key string) (*Result, error) {
	return request(vt.URL("ip_addresses/%s", ip), key)
}

func CheckUrl(url, key string) (*Result, error) {
	return request(vt.URL("urls/%s", url), key)
}

func CheckDomain(url, key string) (*Result, error) {
	return request(vt.URL("domains/%s", url), key)
}

func CheckFileHash(sum, key string) (*Result, error) {
	return request(vt.URL("files/%s", sum), key)
}

func parseEngines(obj *vt.Object, res *Result) {
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

		res.Entries = append(res.Entries, Entry{
			Alert:  slices.Contains(alerts, v["category"].(string)),
			Engine: v["engine_name"].(string),
			Result: v["result"].(string),
		})
	}
}

func parseVerdict(obj *vt.Object, res *Result) {
	res.Bad = countStats(obj, alerts)
	res.All = countStats(obj, []string{
		"malicious",
		"suspicious",
		"undetected",
		"harmless",
		"timeout",
		"confirmed-timeout",
		"failure",
		"type-unsupported",
	})

	res.Alert = res.Bad > 0
	res.Label, _ = obj.GetString("popular_threat_classification.suggested_threat_label")

	if len(res.Label) == 0 {
		switch {
		case res.Bad > 0:
			res.Label = Indecisive
		case res.All > 0:
			res.Label = Clean
		default:
			res.Label = Unrated
		}
	}
}

func countStats(obj *vt.Object, lst []string) (n int) {
	for _, k := range lst {
		v, _ := obj.GetInt64(fmt.Sprintf("last_analysis_stats.%s", k))
		n += int(v)
	}
	return
}

func request(url *url.URL, key string) (*Result, error) {
	res := new(Result)

	api := vt.NewClient(key, vt.WithHTTPClient(client.Default()))

	obj, err := api.GetObject(url)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			res.Label = Unknown
			return res, nil
		}

		return nil, err
	}

	if Verbose > 2 {
		trace(obj)
	}

	parseEngines(obj, res)
	parseVerdict(obj, res)

	if Verbose > 1 {
		for _, e := range res.Entries {
			log.Printf(`result is "%s" by %s`, e.Result, e.Engine)
		}
	}

	return res, nil
}

func trace(obj *vt.Object) {
	buf := bytes.NewBuffer(nil)

	b, err := obj.MarshalJSON()

	if err != nil {
		return
	}

	err = json.Indent(buf, b, "", "  ")

	if err != nil {
		return
	}

	log.Printf("received response:\n%s\n", buf.String())
}
