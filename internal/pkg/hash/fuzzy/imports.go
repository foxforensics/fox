package fuzzy

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/saferwall/pe"

	intern "github.com/cuhsat/fox/v4/internal/pkg/file/convert/bin/pe"
)

var ErrNotSupported = errors.New("file type not supported")

func Lookup(dll string, ord uint32) string {
	// lookup library imports
	v := pe.OrdLookup(dll, uint64(ord), false)

	if len(v) > 0 {
		return v
	}

	// lookup additional imports
	if m, ok := imports[dll]; ok {
		if v, ok = m[ord]; ok {
			return v
		}
	}

	return fmt.Sprintf("ord%d", ord)
}

func GetImports(b []byte, sort bool) ([]string, error) {
	var imp []string

	if !intern.Detect(b) {
		return imp, ErrNotSupported
	}

	p, err := pe.NewBytes(b, &pe.Options{
		DisableCertValidation:      true,
		DisableSignatureValidation: true,
		OmitExportDirectory:        true,
		OmitExceptionDirectory:     true,
		OmitResourceDirectory:      true,
		OmitSecurityDirectory:      true,
		OmitRelocDirectory:         true,
		OmitDebugDirectory:         true,
		OmitArchitectureDirectory:  true,
		OmitGlobalPtrDirectory:     true,
		OmitTLSDirectory:           true,
		OmitLoadConfigDirectory:    true,
		OmitBoundImportDirectory:   true,
		OmitCLRHeaderDirectory:     true,
		OmitCLRMetadata:            true,
	})

	if err != nil {
		return imp, err
	}

	defer func(f *pe.File) {
		_ = f.Close()
	}(p)

	err = p.Parse()

	if err != nil {
		return imp, err
	}

	rep := strings.NewReplacer(".dll", "", ".ocx", "", ".sys", "")

	for _, i := range p.Imports {
		buf := make([]string, 0, len(i.Functions))

		dll := rep.Replace(strings.ToLower(i.Name))

		for _, f := range i.Functions {
			n := strings.ToLower(f.Name)

			if len(f.Name) == 0 {
				n = Lookup(i.Name, f.Ordinal)
			}

			buf = append(buf, fmt.Sprintf("%s.%s", dll, n))
		}

		if sort {
			slices.Sort(buf)
		}

		imp = append(imp, buf...)
	}

	return imp, nil
}
