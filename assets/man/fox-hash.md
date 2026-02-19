% FOX HASH(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hash** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Shows file hashes and checksums. Results will be grouped by path, if more than one _algorithm_ is specified. This mode enforces the **--no-convert** flag.

FLAGS
=====

**-A, --algo**=_name_,...

:   Uses algorithm(s) (default: **SHA256**).

**-a, --all**

:   Uses all algorithms.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ALGORITHMS
==========

Cryptographic hashes (BLAKE family)

:   BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512

Cryptographic hashes (GOST family)

:   GOST2012-256, GOST2012-512

Cryptographic hashes (SHA family)

:   SHA1, SHA224, SHA256, SHA512, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Cryptographic hashes (SKEIN family)

:   SKEIN-224, SKEIN-256, SKEIN-384, SKEIN-512

Cryptographic hashes (MD family)

:   MD2, MD4, MD5, MD6

Cryptographic hashes (other)

:   HAS-160, LSH-256, LSH-512, RIPEMD-160, SHAKE128, SHAKE256, SM3, WHIRLPOOL

Performance hashes

:   DJB2, FNV-1, FNV-1A, MURMUR3, RAPIDHASH, SIPHASH, XXH32, XXH64, XXH3

Similarity hashes

:   IMPFUZZY, IMPHASH, IMPHASH0, SSDEEP, TLSH

Windows specific

:   LM, NT, PE Checksum

Image specific

:   AHASH, DHASH, PHASH

Checksums

:   ADLER32, FLETCHER4, CRC16-CCITT, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

EXAMPLES
========

fox hash -Amd5 files.7z

:   Hash archive contents.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://github.com/cuhsat/fox/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxhunt.dev_>

SEE ALSO
========

**fox(1)**
