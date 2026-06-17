package lookup

import (
	"fmt"
	"log/slog"
	"strings"

	"go.foxforensics.eu/checker/services"
	"go.foxforensics.eu/checker/services/vt"
	"go.foxforensics.eu/hasher/hash"

	"go.foxforensics.eu/fox/v4/internal/pkg/types/carver"
)

func Lookup(a any, verbose int) (bool, error) {
	var res *services.Result
	var err error

	switch v := a.(type) {
	case *carver.String:
		res, err = checkString(v)
	case []byte:
		res, err = checkBytes(v)
	}

	if err != nil {
		return false, err
	}

	if res != nil {
		switch {
		case verbose > 2:
			slog.Info(fmt.Sprintf("lookup:\n%s\n", res.ToJSON()))
		case verbose > 1:
			slog.Info(fmt.Sprintf("lookup: %s [%d/%d]", res.Verdict, res.Stats.Bad, res.Stats.All))
		case verbose > 0:
			slog.Info(fmt.Sprintf("lookup: %s", res.Verdict))
		}
		return res.Stats.Bad > 0, nil
	}

	return false, nil
}

func checkString(s *carver.String) (*services.Result, error) {
	switch {
	case strings.Contains(s.Classes, "IPv6"):
		return vt.CheckIp(s.Value)
	case strings.Contains(s.Classes, "IPv4"):
		return vt.CheckIp(s.Value)
	case strings.Contains(s.Classes, "DNS"):
		return vt.CheckDns(s.Value)
	case strings.Contains(s.Classes, "URL"):
		return vt.CheckUrl(s.Value)
	default:
		return nil, nil
	}
}

func checkBytes(b []byte) (*services.Result, error) {
	return vt.CheckHash(hash.MustSum(hash.SHA256, b))
}
