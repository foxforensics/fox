package virus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"maps"
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

func Test(sum, key string) ([]Entry, error) {
	var e []Entry

	obj, err := api(key).GetObject(vt.URL("files/%s", sum))

	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return e, nil
		}

		return nil, err
	}

	err = trace(obj)

	if err != nil {
		log.Println(err)
	}

	res, err := obj.Get("last_analysis_results")

	if err != nil {
		return nil, err
	}

	m := res.(map[string]any)

	for _, k := range slices.Sorted(maps.Keys(m)) {
		v := m[k].(map[string]any)

		if v["result"] == nil {
			continue
		}

		e = append(e, Entry{
			Alert:  slices.Contains(alerts, v["category"].(string)),
			Engine: v["engine_name"].(string),
			Result: v["result"].(string),
		})
	}

	return e, nil
}

func api(key string) *vt.Client {
	return vt.NewClient(key, vt.WithHTTPClient(client.Default))
}

func trace(obj *vt.Object) error {
	var buf bytes.Buffer

	if !Trace {
		return nil
	}

	b, err := obj.MarshalJSON()

	if err != nil {
		return err
	}

	err = json.Indent(&buf, b, "", "  ")

	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("received response:\n%s", buf.String()))

	return nil
}
