package types

const (
	ADLER32    = "adler32"
	BLAKE2S256 = "blake2s-256"
	BLAKE2B256 = "blake2b-256"
	BLAKE2B384 = "blake2b-384"
	BLAKE2B512 = "blake2b-512"
	BLAKE3256  = "blake3-256"
	BLAKE3512  = "blake3-512"
	CRC32C     = "crc32-c"
	CRC32IEEE  = "crc32-ieee"
	CRC64ECMA  = "crc64-ecma"
	CRC64ISO   = "crc64-iso"
	FNV1       = "fnv-1"
	FNV1A      = "fnv-1a"
	PE         = "pe"
	LM         = "lm"
	NT         = "nt"
	MD2        = "md2"
	MD4        = "md4"
	MD5        = "md5"
	MD6        = "md6"
	MURMUR3    = "murmur3"
	SHA1       = "sha1"
	SHA256     = "sha256"
	SHA512     = "sha512"
	SHA3       = "sha3"
	SHA3224    = "sha3-224"
	SHA3256    = "sha3-256"
	SHA3384    = "sha3-384"
	SHA3512    = "sha3-512"
	SSDEEP     = "ssdeep"
	TLSH       = "tlsh"
	XXH64      = "xxh64"
	XXH3       = "xxh3"
)

const (
	Canonical = "c"
	Hexdump   = "hd"
	Xxd       = "xxd"
	Raw       = "raw"
)

type Event int

const (
	Eventlog Event = iota
	Journal
)

type Heap int

const (
	Stdin Heap = iota
	Stdout
	Stderr
	Regular
	Deflate
	Defined
)
