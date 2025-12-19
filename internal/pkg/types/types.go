package types

const Database = "fox.db"

const Buffer = 4095

const (
	ADLER32   = "adler32"
	CRC32C    = "crc32-c"
	CRC32IEEE = "crc32-ieee"
	CRC64ECMA = "crc64-ecma"
	CRC64ISO  = "crc64-iso"
	PE        = "pe"
	LM        = "lm"
	NT        = "nt"
	MD2       = "md2"
	MD4       = "md4"
	MD5       = "md5"
	MD6       = "md6"
	SHA1      = "sha1"
	SHA256    = "sha256"
	SHA3      = "sha3"
	SHA3224   = "sha3-224"
	SHA3256   = "sha3-256"
	SHA3384   = "sha3-384"
	SHA3512   = "sha3-512"
	SSDEEP    = "ssdeep"
	TLSH      = "tlsh"
	XXH64     = "xxh64"
	XXH3      = "xxh3"
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
	String
	Regular
	Deflate
)
