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
	Clean      = "clean"
	Indecisive = "indecisive"
	Unrated    = "unrated"
)

// Trace API responses
var Trace bool

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
}

func TestIp(ip, key string) (*Result, error) {
	return request(vt.URL("ip_addresses/%s", ip), key)
}

func TestUrl(url, key string) (*Result, error) {
	return request(vt.URL("urls/%s", url), key)
}

func TestDomain(url, key string) (*Result, error) {
	return request(vt.URL("domains/%s", url), key)
}

func TestFileHash(sum, key string) (*Result, error) {
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
	bad := countStats(obj, alerts)
	all := countStats(obj, []string{
		"malicious",
		"suspicious",
		"undetected",
		"harmless",
		"timeout",
		"confirmed-timeout",
		"failure",
		"type-unsupported",
	})

	res.Alert = bad > 0
	res.Label, _ = obj.GetString("popular_threat_classification.suggested_threat_label")

	if len(res.Label) == 0 {
		switch {
		case bad > 0:
			res.Label = Indecisive
		case all > 0:
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
			res.Label = "unknown"
			return res, nil
		}

		return nil, err
	}

	err = trace(obj)

	if err != nil {
		log.Println(err)
	}

	parseEngines(obj, res)
	parseVerdict(obj, res)

	return res, nil
}

func trace(obj *vt.Object) error {
	buf := bytes.NewBuffer(nil)

	if !Trace {
		return nil
	}

	b, err := obj.MarshalJSON()

	if err != nil {
		return err
	}

	err = json.Indent(buf, b, "", "  ")

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("received response:\n%s", buf.String()))

	return nil
}
