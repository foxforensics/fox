package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/glaslos/ssdeep"
	"github.com/glaslos/tlsh"
	"github.com/htruong/go-md2"
	"github.com/zeebo/xxh3"
	"golang.org/x/crypto/md4"

	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/lm"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/nt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

func Secure(t string) bool {
	switch strings.ToLower(t) {
	case types.SHA256:
		fallthrough
	case types.SHA3:
		fallthrough
	case types.SHA3224, types.SHA3256, types.SHA3384, types.SHA3512:
		return true
	default:
		return false
	}
}

func Sum(t string, b []byte) (string, error) {
	var imp hash.Hash

	switch strings.ToLower(t) {
	case types.CRC32C:
		imp = crc32.New(crc32.MakeTable(crc32.Castagnoli))
	case types.CRC32IEEE:
		imp = crc32.NewIEEE()
	case types.CRC64ECMA:
		imp = crc64.New(crc64.MakeTable(crc64.ECMA))
	case types.CRC64ISO:
		imp = crc64.New(crc64.MakeTable(crc64.ISO))
	case types.PE:
		imp = pe.New()
	case types.LM:
		imp = lm.New()
	case types.NT:
		imp = nt.New()
	case types.MD2:
		imp = md2.New()
	case types.MD4:
		imp = md4.New()
	case types.MD5:
		imp = md5.New()
	case types.SHA1:
		imp = sha1.New()
	case types.SHA256:
		imp = sha256.New()
	case types.SHA3:
		fallthrough
	case types.SHA3224:
		imp = sha3.New224()
	case types.SHA3256:
		imp = sha3.New256()
	case types.SHA3384:
		imp = sha3.New384()
	case types.SHA3512:
		imp = sha3.New512()
	case types.SSDEEP:
		imp = ssdeep.New()
	case types.TLSH:
		imp = tlsh.New()
	case types.XXH3:
		imp = xxh3.New()
	case types.XXH64:
		imp = xxhash.New()
	default:
		return "", errors.New("algorithm not recognized")
	}

	imp.Reset()

	if _, err := imp.Write(b); err != nil {
		return "", err
	}

	b = imp.Sum(nil)

	if len(b) == 0 {
		return "", errors.New("input size to small")
	}

	switch t {
	case types.SSDEEP:
		return fmt.Sprintf("%s", b), nil
	case types.TLSH:
		return fmt.Sprintf("T1%x", b), nil
	default:
		return fmt.Sprintf("%x", b), nil
	}
}
