package vt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/VirusTotal/vt-go"

	"github.com/cuhsat/fox/v4/internal/pkg/text"
	"github.com/cuhsat/fox/v4/internal/pkg/types/client"
)

func GetLabel(sum, key string) (string, error) {
	obj, err := api(key).GetObject(files(sum))

	if err != nil {
		return "", err
	}

	// malicious count
	mal := obj.MustGetInt64("last_analysis_stats.malicious")
	mal += obj.MustGetInt64("last_analysis_stats.suspicious")

	// overall count
	all := obj.MustGetInt64("last_analysis_stats.undetected")
	all += obj.MustGetInt64("last_analysis_stats.harmless")
	all += obj.MustGetInt64("last_analysis_stats.timeout")
	all += obj.MustGetInt64("last_analysis_stats.confirmed-timeout")
	all += obj.MustGetInt64("last_analysis_stats.failure")
	all += obj.MustGetInt64("last_analysis_stats.type-unsupported")

	// popular label
	lbl := obj.MustGetString("popular_threat_classification.suggested_threat_label")

	return fmt.Sprintf("%d of %d %s", mal, mal+all, lbl), nil
}

func GetResult(sum, key string) (string, error) {
	obj, err := api(key).GetObject(files(sum))

	if err != nil {
		return "", err
	}

	res, err := obj.Get("last_analysis_results")

	var ar = res.(map[string]any)
	var sb strings.Builder

	for _, k := range slices.Sorted(maps.Keys(ar)) {
		v := ar[k].(map[string]any)["result"]

		if v == nil {
			v = "-" // unknown
		}

		sb.WriteString(fmt.Sprintf("%30s %s\n", text.Hide(k), v))
	}

	return strings.Trim(sb.String(), "\n"), nil
}

func GetReport(sum, key string, flat bool) (string, error) {
	res, err := api(key).Get(files(sum))

	if err != nil {
		return "", err
	}

	if flat {
		return string(res.Data), nil
	}

	buf := bytes.NewBuffer(nil)

	_ = json.Indent(buf, res.Data, "", "  ")

	return buf.String(), nil
}

func api(key string) *vt.Client {
	return vt.NewClient(key, vt.WithHTTPClient(client.Default))
}

func files(sum string) *url.URL {
	return vt.URL("files/%s", sum)
}
