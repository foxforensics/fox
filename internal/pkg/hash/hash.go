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

	"github.com/cuhsat/fox/v4/internal/pkg/hash/crc/kermit"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/crypto/blake3"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/crypto/shake"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/fast/xxh"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/fuzzy/impfuzzy"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/fuzzy/imphash"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/image/image"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/lm"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/nt"
	"github.com/cuhsat/fox/v4/internal/pkg/hash/windows/pe"
	"github.com/cuhsat/fox/v4/internal/pkg/types"
)

var Algorithms = []struct {
	Name   string
	Secure bool
}{
	{types.ADLER32, false},
	{types.AVERAGE, false},
	{types.BLAKE2S256, true},
	{types.BLAKE2B256, true},
	{types.BLAKE2B384, true},
	{types.BLAKE2B512, true},
	{types.BLAKE3256, true},
	{types.BLAKE3512, true},
	{types.BLOCKMEAN, false},
	{types.CRC16CCITT, false},
	{types.CRC32C, false},
	{types.CRC32IEEE, false},
	{types.CRC64ECMA, false},
	{types.CRC64ISO, false},
	{types.DIFFERENCE, false},
	{types.FLETCHER4, false},
	{types.FNV1, false},
	{types.FNV1A, false},
	{types.GOST2012256, true},
	{types.GOST2012512, true},
	{types.HAS160, false},
	{types.IMPFUZZY, false},
	{types.IMPHASH, false},
	{types.IMPHASH0, false},
	{types.LM, false},
	{types.LSH256, true},
	{types.LSH512, true},
	{types.MARRHILDRETH, false},
	{types.MD2, false},
	{types.MD4, false},
	{types.MD5, false},
	{types.MD6, false},
	{types.MEDIAN, false},
	{types.MURMUR3, false},
	{types.NT, false},
	{types.PDQ, false},
	{types.PE, false},
	{types.PHASH, false},
	{types.RAPIDHASH, false},
	{types.RASH, false},
	{types.RIPEMD160, true},
	{types.SHAKE128, true},
	{types.SHAKE256, true},
	{types.SHA1, false},
	{types.SHA256, true},
	{types.SHA512, true},
	{types.SHA3, true},
	{types.SHA3224, true},
	{types.SHA3256, true},
	{types.SHA3384, true},
	{types.SHA3512, true},
	{types.SIPHASH, false},
	{types.SKEIN224, true},
	{types.SKEIN256, true},
	{types.SKEIN384, true},
	{types.SKEIN512, true},
	{types.SM3, true},
	{types.SSDEEP, false},
	{types.STREEBOG256, true},
	{types.STREEBOG512, true},
	{types.TLSH, false},
	{types.WHASH, false},
	{types.WHIRLPOOL, true},
	{types.XXH3, false},
	{types.XXH32, false},
	{types.XXH64, false},
}

func IsSecure(algo string) bool {
	for _, a := range Algorithms {
		if a.Name == strings.ToLower(algo) {
			return a.Secure
		}
	}

	return false
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
	case types.AVERAGE:
		imp = image.New(image.Average)
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
	case types.BLOCKMEAN:
		imp = image.New(image.BlockMean)
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
	case types.DIFFERENCE:
		imp = image.New(image.Difference)
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
	case types.MARRHILDRETH:
		imp = image.New(image.MarrHildreth)
	case types.MD2:
		imp = md2.New()
	case types.MD4:
		imp = md4.New()
	case types.MD5:
		imp = md5.New()
	case types.MD6:
		imp = md6.New256()
	case types.MEDIAN:
		imp = image.New(image.Median)
	case types.MURMUR3:
		imp = murmur3.New64() // Murmur3f
	case types.NT:
		imp = nt.New()
	case types.PDQ:
		imp = image.New(image.PDQ)
	case types.PE:
		imp = pe.New()
	case types.PHASH:
		imp = image.New(image.PHash)
	case types.RAPIDHASH:
		imp = rapidhash.New()
	case types.RASH:
		imp = image.New(image.RASH)
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
	case types.WHASH:
		imp = image.New(image.WHash)
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
