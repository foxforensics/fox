package pe

import (
	"encoding/json"
	"log"

	"github.com/saferwall/pe"

	"github.com/cuhsat/fox/v4/internal/pkg/data"
)

const Magic = "MZ"

func Detect(b []byte) bool {
	return data.HasMagic(b, 0, []byte(Magic))
}

func Convert(b []byte) ([]byte, error) {
	p, err := pe.NewBytes(b, &pe.Options{
		SectionEntropy:         true,
		OmitExceptionDirectory: true, // cut for clarity
		OmitRelocDirectory:     true, // cut for clarity
	})

	if err != nil {
		return nil, err
	}

	err = p.Parse()

	if err != nil {
		return nil, err
	}

	for _, a := range p.Anomalies {
		log.Printf("warning: %s!\n", a)
	}

	return json.Marshal(p)
}
