% FOX HASH(1) Version 5 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hash** \[_flags_ ...] \[**list** | _paths_ ...]

DESCRIPTION
===========

Show different file hashes and checksums. Results will be grouped by path, if more than one _algorithm_ is specified.

FLAGS
=====

**-H, --hash**=_name_,...

:   Use hash algorithm(s) (default: **SHA256**).

**-a, --all**

:   Show all hashes and checksums.

**-j, --json**

:   Show results as JSON objects.

**-J, --jsonl**

:   Show results as JSON lines.

Filter Flags
------------
All filter flags can only be used while using a single hash algorithm. 

**-B, --include**=_file_

:   Include only known bad hashes, which are loaded from the given _file_.

**-G, --exclude**=_file_

:   Exclude all known good hashes, which are loaded from the given _file_. For a list of known good hashes, visit the **National Software Reference Library** at <_https://www.nist.gov/itl/csd/secure-systems-and-applications/national-software-reference-library-nsrl_>.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. If **list** is specified as _path_, only the list of the built-in algorithms will be shown. To refer to paths inside archives, use the archive::file notation.

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

:   FNV-1, FNV-1A, MURMUR3, RAPIDHASH, SIPHASH, XXH32, XXH64, XXH3

Perceptual hashes

:   AVERAGE, DIFFERENCE, MEDIAN, PHASH, WHASH, MARRHILDRETH, BLOCKMEAN, PDQ, RASH

Similarity hashes

:   IMPFUZZY, IMPHASHO, IMPHASHS, SDHASH, SSDEEP, TLSH

Windows specific

:   LM, NT, PE

Unix specific

:   BSD, ELF, SYSV

Checksums

:   ADLER32, FLETCHER4, LUHN, CRC16-CCITT, CRC32-C, CRC32-K, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

EXAMPLES
========

$ fox hash -Hmd5 files.7z

:   Hash archive contents as MD5.

$ fox hash -Himpfuzzy *.exe

:   Hash binaries for similarity.

$ fox hash -Pinfected ioc.zip::ioc.exe

:   Hash binary inside an archive.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://foxforensics.eu/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxforensics.eu_>

SEE ALSO
========

**fox(1)**, **sum(1)**, **md5sum(1)**, **sha1sum(1)**, **sha256sum(1)**
