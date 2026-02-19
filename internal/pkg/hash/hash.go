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
	"github.com/cuhsat/go-hash/skein"
	"github.com/cuhsat/go-hash/streebog"
	"github.com/cuhsat/go-krypto/has160"
	"github.com/cuhsat/go-krypto/lsh256"
	"github.com/cuhsat/go-krypto/lsh512"
	"github.com/dchest/siphash"
	"github.com/glaslos/ssdeep"
	"github.com/glaslos/tlsh"
	"github.com/htruong/go-md2"
	"github.com/jzelinskie/whirlpool"
	"github.com/pedroalbanese/md6"
	"github.com/spaolacci/murmur3"
	"github.com/tjfoc/gmsm/v2/sm3"
	"github.com/zeebo/xxh3"
	"go.dw1.io/rapidhash"
	"go.solidsystem.no/fletcher4"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"

	"github.com/cuhsat/fox/v4/internal/pkg/hash/crypto/blake3"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/crypto/shake"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/fast/djb2"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/fast/xxh"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/fuzzy/impfuzzy"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/fuzzy/imphash"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/other/image"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/other/kermit"
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
	types.CRC16CCITT,
	types.CRC32C,
	types.CRC32IEEE,
	types.CRC64ECMA,
	types.CRC64ISO,
	types.DHASH,
	types.DJB2,
	types.FLETCHER4,
	types.FNV1,
	types.FNV1A,
	types.GOST2012256,
	types.GOST2012512,
	types.HAS160,
	types.IMPFUZZY,
	types.IMPHASH,
	types.IMPHASH0,
	types.LSH256,
	types.LSH512,
	types.LM,
	types.MD2,
	types.MD4,
	types.MD5,
	types.MD6,
	types.MURMUR3,
	types.NT,
	types.PE,
	types.PHASH,
	types.RAPIDHASH,
	types.RIPEMD160,
	types.SHAKE128,
	types.SHAKE256,
	types.SHA1,
	types.SHA256,
	types.SHA512,
	types.SHA3,
	types.SHA3224,
	types.SHA3256,
	types.SHA3384,
	types.SHA3512,
	types.SIPHASH,
	types.SKEIN224,
	types.SKEIN256,
	types.SKEIN384,
	types.SKEIN512,
	types.SM3,
	types.SSDEEP,
	types.TLSH,
	types.WHIRLPOOL,
	types.XXH3,
	types.XXH32,
	types.XXH64,
}

var secure = []string{
	types.BLAKE2S256,
	types.BLAKE2B256,
	types.BLAKE2B384,
	types.BLAKE2B512,
	types.BLAKE3256,
	types.BLAKE3512,
	types.GOST2012256,
	types.GOST2012512,
	types.LSH256,
	types.LSH512,
	types.RIPEMD160,
	types.SHAKE128,
	types.SHAKE256,
	types.SHA256,
	types.SHA512,
	types.SHA3,
	types.SHA3224,
	types.SHA3256,
	types.SHA3384,
	types.SHA3512,
	types.SKEIN224,
	types.SKEIN256,
	types.SKEIN384,
	types.SKEIN512,
	types.SM3,
	types.STREEBOG256,
	types.STREEBOG512,
	types.WHIRLPOOL,
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

	ssdeep.Force = true

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
	case types.CRC16CCITT:
		imp = kermit.New()
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
	case types.DJB2:
		imp = djb2.New()
	case types.FLETCHER4:
		imp = fletcher4.New()
	case types.FNV1:
		imp = fnv.New128()
	case types.FNV1A:
		imp = fnv.New128a()
	case types.GOST2012256, types.STREEBOG256:
		imp = streebog.New256()
	case types.GOST2012512, types.STREEBOG512:
		imp = streebog.New512()
	case types.HAS160:
		imp = has160.New()
	case types.IMPFUZZY:
		imp = impfuzzy.New()
	case types.IMPHASH:
		imp = imphash.New()
	case types.IMPHASH0:
		imp = imphash.NewStable()
	case types.LM:
		imp = lm.New()
	case types.LSH256:
		imp = lsh256.New()
	case types.LSH512:
		imp = lsh512.New()
	case types.MD2:
		imp = md2.New()
	case types.MD4:
		imp = md4.New()
	case types.MD5:
		imp = md5.New()
	case types.MD6:
		imp = md6.New256()
	case types.MURMUR3:
		imp = murmur3.New64() // Murmur3f
	case types.NT:
		imp = nt.New()
	case types.PE:
		imp = pe.New()
	case types.PHASH:
		imp = image.NewPHash()
	case types.RAPIDHASH:
		imp = rapidhash.New()
	case types.RIPEMD160:
		imp = ripemd160.New()
	case types.SHA1:
		imp = sha1.New()
	case types.SHA224:
		imp = sha256.New224()
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
	case types.SHAKE128:
		imp = shake.New128()
	case types.SHAKE256:
		imp = shake.New256()
	case types.SIPHASH:
		imp = siphash.New(make([]byte, 16)) // SipHash-2-4 with zero key
	case types.SKEIN224:
		imp = skein.NewHash224()
	case types.SKEIN256:
		imp = skein.NewHash256()
	case types.SKEIN384:
		imp = skein.NewHash384()
	case types.SKEIN512:
		imp = skein.NewHash512()
	case types.SM3:
		imp = sm3.New()
	case types.SSDEEP:
		imp = ssdeep.New()
	case types.TLSH:
		imp = tlsh.New()
	case types.WHIRLPOOL:
		imp = whirlpool.New()
	case types.XXH3:
		imp = xxh3.New()
	case types.XXH32:
		imp = xxh.New()
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
	case types.SSDEEP, types.IMPFUZZY:
		return fmt.Sprintf("%s", data), nil
	case types.TLSH:
		return fmt.Sprintf("T1%x", data), nil
	default:
		return fmt.Sprintf("%x", data), nil
	}
}
