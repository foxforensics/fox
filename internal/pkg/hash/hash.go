package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"log"
	"slices"
	"strings"

	"github.com/cespare/xxhash"
	"github.com/glaslos/ssdeep"
	"github.com/glaslos/tlsh"
	"github.com/htruong/go-md2"
	"github.com/pedroalbanese/md6"
	"github.com/spaolacci/murmur3"
	"github.com/zeebo/xxh3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"

	"github.com/cuhsat/fox/v4/internal/pkg/hash/crypto/blake3"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/others/image"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/lm"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/nt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Algorithms = []string{
	types.ADLER32,
	types.AHASH,
	types.BLAKE2S256,
	types.BLAKE2B256,
	types.BLAKE2B384,
	types.BLAKE2B512,
	types.BLAKE3256,
	types.BLAKE3512,
	types.CRC32C,
	types.CRC32IEEE,
	types.CRC64ECMA,
	types.CRC64ISO,
	types.DHASH,
	types.FNV1,
	types.FNV1A,
	types.LM,
	types.MD2,
	types.MD4,
	types.MD5,
	types.MD6,
	types.MURMUR3,
	types.NT,
	types.PE,
	types.PHASH,
	types.SHA1,
	types.SHA256,
	types.SHA512,
	types.SHA3,
	types.SHA3224,
	types.SHA3256,
	types.SHA3384,
	types.SHA3512,
	types.SSDEEP,
	types.TLSH,
	types.XXH64,
	types.XXH3,
}

var secure = []string{
	types.BLAKE2S256,
	types.BLAKE2B256,
	types.BLAKE2B384,
	types.BLAKE2B512,
	types.BLAKE3256,
	types.BLAKE3512,
	types.SHA256,
	types.SHA512,
	types.SHA3,
	types.SHA3224,
	types.SHA3256,
	types.SHA3384,
	types.SHA3512,
}

func IsSecure(algo string) bool {
	return slices.Contains(secure, strings.ToLower(algo))
}

func MustSum(algo string, data []byte) string {
	sum, err := Sum(algo, data)

	if err != nil {
		log.Fatalln(err)
	}

	return sum
}

func Sum(algo string, data []byte) (string, error) {
	var imp hash.Hash

	switch strings.ToLower(algo) {
	case types.ADLER32:
		imp = adler32.New()
	case types.AHASH:
		imp = image.NewAHash()
	case types.BLAKE2B256:
		imp, _ = blake2b.New256(nil)
	case types.BLAKE2B384:
		imp, _ = blake2b.New384(nil)
	case types.BLAKE2B512:
		imp, _ = blake2b.New512(nil)
	case types.BLAKE2S256:
		imp, _ = blake2s.New256(nil)
	case types.BLAKE3256:
		imp = blake3.New256()
	case types.BLAKE3512:
		imp = blake3.New512()
	case types.CRC32C:
		imp = crc32.New(crc32.MakeTable(crc32.Castagnoli))
	case types.CRC32IEEE:
		imp = crc32.NewIEEE()
	case types.CRC64ECMA:
		imp = crc64.New(crc64.MakeTable(crc64.ECMA))
	case types.CRC64ISO:
		imp = crc64.New(crc64.MakeTable(crc64.ISO))
	case types.DHASH:
		imp = image.NewDHash()
	case types.FNV1:
		imp = fnv.New128()
	case types.FNV1A:
		imp = fnv.New128a()
	case types.LM:
		imp = lm.New()
	case types.MD2:
		imp = md2.New()
	case types.MD4:
		imp = md4.New()
	case types.MD5:
		imp = md5.New()
	case types.MD6:
		imp = md6.New256()
	case types.MURMUR3:
		imp = murmur3.New64() // MURMUR3F
	case types.NT:
		imp = nt.New()
	case types.PE:
		imp = pe.New()
	case types.PHASH:
		imp = image.NewPHash()
	case types.SHA1:
		imp = sha1.New()
	case types.SHA256:
		imp = sha256.New()
	case types.SHA512:
		imp = sha512.New()
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

	if _, err := imp.Write(data); err != nil {
		return "", err
	}

	data = imp.Sum(nil)

	if len(data) == 0 {
		return "", errors.New("input size to small")
	}

	switch algo {
	case types.SSDEEP:
		return fmt.Sprintf("%s", data), nil
	case types.TLSH:
		return fmt.Sprintf("T1%x", data), nil
	default:
		return fmt.Sprintf("%x", data), nil
	}
}
