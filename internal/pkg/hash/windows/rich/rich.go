package rich

import (
	"crypto/md5"
	"encoding/hex"
	"hash"

	"github.com/saferwall/pe"
)

type Rich struct {
	buf []byte
}

func New() hash.Hash {
	return new(Rich)
}

func (h *Rich) BlockSize() int {
	return md5.Size
}

func (h *Rich) Size() int {
	return md5.Size
}

func (h *Rich) Reset() {
	h.buf = h.buf[:0]
}

func (h *Rich) Write(b []byte) (n int, err error) {
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
		return 0, err
	}

	h.buf, err = hex.DecodeString(p.RichHeaderHash())

	_ = p.Close()

	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (h *Rich) Sum(_ []byte) []byte {
	return h.buf
}
