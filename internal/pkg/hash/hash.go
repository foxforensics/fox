package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"errors"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/cxmcc/tiger"
	"github.com/dchest/siphash"
	"github.com/glaslos/ssdeep"
	"github.com/glaslos/tlsh"
	"github.com/htruong/go-md2"
	"github.com/jzelinskie/whirlpool"
	"github.com/zeebo/xxh3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"

	"github.com/cuhsat/fox/v4/internal/pkg/hash/blake3"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/sdhash"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

func Sum(t string, b []byte) ([]byte, error) {
	var imp hash.Hash

	switch strings.ToLower(t) {
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
	case types.SHA3, types.SHA3224:
		imp = sha3.New224()
	case types.SHA3256:
		imp = sha3.New256()
	case types.SHA3384:
		imp = sha3.New384()
	case types.SHA3512:
		imp = sha3.New512()
	case types.BLAKE2S256:
		imp, _ = blake2s.New256(nil)
	case types.BLAKE2B256:
		imp, _ = blake2b.New256(nil)
	case types.BLAKE2B384:
		imp, _ = blake2b.New384(nil)
	case types.BLAKE2B512:
		imp, _ = blake2b.New512(nil)
	case types.BLAKE3256:
		imp = blake3.New256()
	case types.BLAKE3512:
		imp = blake3.New512()
	case types.RIPEMD160:
		imp = ripemd160.New()
	case types.TIGER:
		imp = tiger.New()
	case types.TIGER2:
		imp = tiger.New2()
	case types.WHIRLPOOL:
		imp = whirlpool.New()
	case types.FNV1:
		imp = fnv.New64()
	case types.FNV1A:
		imp = fnv.New64a()
	case types.SIPHASH64:
		imp = siphash.New(nil)
	case types.SIPHASH128:
		imp = siphash.New128(nil)
	case types.XXH64:
		imp = xxhash.New()
	case types.XXH3:
		imp = xxh3.New()
	case types.SDHASH:
		imp = sdhash.New()
	case types.SSDEEP:
		imp = ssdeep.New()
	case types.TLSH:
		imp = tlsh.New()
	case types.ADLER32:
		imp = adler32.New()
	case types.CRC32IEEE:
		imp = crc32.NewIEEE()
	case types.CRC64ISO:
		imp = crc64.New(crc64.MakeTable(crc64.ISO))
	case types.CRC64ECMA:
		imp = crc64.New(crc64.MakeTable(crc64.ECMA))
	default:
		return nil, errors.New("algorithm not recognized")
	}

	imp.Reset()

	if _, err := imp.Write(b); err != nil {
		return nil, err
	}

	return imp.Sum(nil), nil
}
