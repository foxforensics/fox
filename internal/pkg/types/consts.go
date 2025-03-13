package types

const Database = "fox.db"

const (
	MD5       = "md5"
	SHA1      = "sha1"
	SHA256    = "sha256"
	SHA3      = "sha3"
	SHA3224   = "sha3-224"
	SHA3256   = "sha3-256"
	SHA3384   = "sha3-384"
	SHA3512   = "sha3-512"
	BLAKE3256 = "blake3-256"
	BLAKE3512 = "blake3-512"

	FNV1  = "fnv-1"
	FNV1A = "fnv-1a"
	XXH64 = "xxh64"
	XXH3  = "xxh3"

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
