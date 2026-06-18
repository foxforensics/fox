package lookup

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/VirusTotal/vt-go"
	"go.foxforensics.eu/hasher/hash"

	"go.foxforensics.eu/fox/v4/internal/pkg/types/carver"
	"go.foxforensics.eu/fox/v4/internal/pkg/types/client"
)

func Lookup(key string, a any) (bool, error) {
	switch v := a.(type) {
	case []byte:
		return request(key, vt.URL("files/%s", hash.MustSum(hash.SHA256, v)))

	case *carver.String:
		switch {
		case strings.Contains(v.Classes, "IPv6"):
			return request(key, vt.URL("ip_addresses/%s", v.Value))
		case strings.Contains(v.Classes, "IPv4"):
			return request(key, vt.URL("ip_addresses/%s", v.Value))
		case strings.Contains(v.Classes, "DNS"):
			return request(key, vt.URL("domains/%s", v.Value))
		case strings.Contains(v.Classes, "URL"):
			return request(key, vt.URL("urls/%s", v.Value))
		default:
			return false, nil
		}

	default:
		return false, nil
	}
}

func request(key string, url *url.URL) (bool, error) {
	c := vt.NewClient(key, vt.WithHTTPClient(client.Http()))

	obj, err := c.GetObject(url)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "X-Apikey header is missing"):
			return false, errors.New("API key is missing")
		case strings.Contains(err.Error(), "not found"):
			return false, nil
		default:
			return false, err
		}
	}

	if b, err := obj.MarshalJSON(); err == nil {
		slog.Debug(fmt.Sprintf("lookup: %s\n", b))
	} else {
		slog.Error(err.Error())
	}

	for _, k := range []string{
		"last_analysis_stats.malicious",
		"last_analysis_stats.suspicious",
	} {
		v, _ := obj.GetInt64(k)

		// at least one bad stat
		if int(v) > 0 {
			return true, nil
		}
	}

	return false, nil
}
