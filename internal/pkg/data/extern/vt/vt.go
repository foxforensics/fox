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
	Bad     int64
	All     int64
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
	var (
		mal, _ = obj.GetInt64("last_analysis_stats.malicious")
		sus, _ = obj.GetInt64("last_analysis_stats.suspicious")
		un, _  = obj.GetInt64("last_analysis_stats.undetected")
		ha, _  = obj.GetInt64("last_analysis_stats.harmless")
		to, _  = obj.GetInt64("last_analysis_stats.timeout")
		cto, _ = obj.GetInt64("last_analysis_stats.confirmed-timeout")
		fa, _  = obj.GetInt64("last_analysis_stats.failure")
		tu, _  = obj.GetInt64("last_analysis_stats.type-unsupported")
	)

	res.Bad = mal + sus
	res.All = res.Bad + un + ha + to + cto + fa + tu
	res.Label, _ = obj.GetString("popular_threat_classification.suggested_threat_label")
	res.Alert = res.Bad > 0

	if len(res.Label) == 0 {
		switch {
		case res.Bad > 0:
			res.Label = "indecisive"
		case res.All > 0:
			res.Label = "clean"
		default:
			res.Label = "unrated"
		}
	}
}

func request(url *url.URL, key string) (*Result, error) {
	res := new(Result)

	api := vt.NewClient(key, vt.WithHTTPClient(client.Default))

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
