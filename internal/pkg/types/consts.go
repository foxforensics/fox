package types

const Database = "fox.db"

const (
	BLAKE2S256 = "blake2s-256"
	BLAKE2B256 = "blake2b-256"
	BLAKE2B384 = "blake2b-384"
	BLAKE2B512 = "blake2b-512"
	BLAKE3256  = "blake3-256"
	BLAKE3512  = "blake3-512"
	GOST256    = "gost-256"
	GOST512    = "gost-512"
	MD2        = "md2"
	MD4        = "md4"
	MD5        = "md5"
	RIPEMD160  = "ripemd-160"
	SHA1       = "sha1"
	SHA256     = "sha256"
	SHA3       = "sha3"
	SHA3224    = "sha3-224"
	SHA3256    = "sha3-256"
	SHA3384    = "sha3-384"
	SHA3512    = "sha3-512"
	TIGER      = "tiger"
	TIGER2     = "tiger2"
	WHIRLPOOL  = "whirlpool"

	FNV1       = "fnv-1"
	FNV1A      = "fnv-1a"
	SIPHASH64  = "siphash-64"
	SIPHASH128 = "siphash-128"
	XXH64      = "xxh64"
	XXH3       = "xxh3"

	SDHASH = "sdhash"
	SSDEEP = "ssdeep"
	TLSH   = "tlsh"

	ADLER32   = "adler32"
	CRC32IEEE = "crc32-ieee"
	CRC64ECMA = "crc64-ecma"
	CRC64ISO  = "crc64-iso"
)

const (
	Canonical = "c"
	Hexdump   = "hd"
	Xxd       = "xxd"
	Raw       = "raw"
)

const (
	Logstash = "http://localhost:8080"
	Splunk   = "http://localhost:8088/services/collector/event/1.0"
)

type Heap int

const (
	Stdin Heap = iota
	Stdout
	Stderr
	Regular
	Deflate
)
