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
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/lm"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/nt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Algorithms = []string{
	"adler32",
	"blake2s-256",
	"blake2b-256",
	"blake2b-384",
	"blake2b-512",
	"blake3-256",
	"blake3-512",
	"crc32-c",
	"crc32-ieee",
	"crc64-ecma",
	"crc64-iso",
	"fnv-1",
	"fnv-1a",
	"pe",
	"lm",
	"nt",
	"md2",
	"md4",
	"md5",
	"md6",
	"murmur3",
	"sha1",
	"sha256",
	"sha512",
	"sha3",
	"sha3-224",
	"sha3-256",
	"sha3-384",
	"sha3-512",
	"ssdeep",
	"tlsh",
	"xxh64",
	"xxh3",
}

var isSecure = []string{
	"blake2s-256",
	"blake2b-256",
	"blake2b-384",
	"blake2b-512",
	"blake3-256",
	"blake3-512",
	"sha256",
	"sha512",
	"sha3",
	"sha3-224",
	"sha3-256",
	"sha3-384",
	"sha3-512",
}

func IsSecure(algo string) bool {
	return slices.Contains(isSecure, strings.ToLower(algo))
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
	case types.FNV1:
		imp = fnv.New128()
	case types.FNV1A:
		imp = fnv.New128a()
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
	case types.MD6:
		imp = md6.New256()
	case types.MURMUR3:
		imp = murmur3.New64() // MURMUR3F
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
