package lookup

import (
	"errors"
	"log"
	"strings"

	"go.foxforensics.dev/checker/services"
	"go.foxforensics.dev/checker/services/vt"
	"go.foxforensics.dev/hasher/hash"

	"go.foxforensics.dev/fox/v4/internal/pkg/types/carver"
)

func Lookup(a any, verbose int) bool {
	var res *services.Result
	var err error

	switch v := a.(type) {
	case *carver.String:
		res, err = checkString(v)
	case []byte:
		res, err = checkBytes(v)
	}

	if errors.Is(err, services.ErrNoApiKey) {
		log.Fatalln("VirusTotal API key not set")
	} else if err != nil {
		log.Fatalln(err)
	}

	if res != nil {
		switch {
		case verbose > 2:
			log.Printf("lookup:\n%s\n", res.ToJSON())
		case verbose > 1:
			log.Printf("lookup: %s [%d/%d]\n", res.Verdict, res.Stats.Bad, res.Stats.All)
		case verbose > 0:
			log.Printf("lookup: %s\n", res.Verdict)
		}
		return res.Stats.Bad > 0
	}

	return false
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
		return nil, errors.ErrUnsupported
	}
}

func checkBytes(b []byte) (*services.Result, error) {
	return vt.CheckHash(hash.MustSum(hash.SHA256, b))
}
