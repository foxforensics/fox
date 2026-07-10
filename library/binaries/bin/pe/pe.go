package pe

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/saferwall/pe"
	"go.foxforensics.eu/fox/v5/library"
)

const Magic = "MZ"

func Detect(b []byte) bool {
	return library.HasMagic(b, 0, []byte(Magic))
}

func Convert(b []byte) ([]byte, error) {
	p, err := pe.NewBytes(b, &pe.Options{
		SectionEntropy:         true,
		OmitExceptionDirectory: true, // cut for clarity
		OmitRelocDirectory:     true, // cut for clarity
	})

	if err != nil {
		return b, err
	}

	defer func() {
		if err = p.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	err = p.Parse()

	if err != nil {
		return b, err
	}

	for _, a := range p.Anomalies {
		slog.Warn(fmt.Sprintf("anomaly: %s!", a))
	}

	return json.Marshal(p)
}
