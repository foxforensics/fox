package pe

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/saferwall/pe"
	"go.foxforensics.eu/fox/v4/internal/pkg"
)

const Magic = "MZ"

func Detect(b []byte) bool {
	return pkg.HasMagic(b, 0, []byte(Magic))
}

func Convert(b []byte) ([]byte, error) {
	p, err := pe.NewBytes(b, &pe.Options{
		SectionEntropy:         true,
		OmitExceptionDirectory: true, // cut for clarity
		OmitRelocDirectory:     true, // cut for clarity
	})

	defer func() {
		_ = p.Close()
	}()

	if err != nil {
		return b, err
	}

	err = p.Parse()

	if err != nil {
		return b, err
	}

	for _, a := range p.Anomalies {
		slog.Warn(fmt.Sprintf("anomaly: %s!", a))
	}

	return json.Marshal(p)
}
