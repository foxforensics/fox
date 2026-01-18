![fox](assets/logo.png "fox")

The Forensic Examiners Swiss Army Knife. Providing many useful features to leverage your forensic examination process. Standalone binaries available for Windows, Linux and macOS.

![Go Report](https://goreportcard.com/badge/github.com/cuhsat/fox/v4?style=for-the-badge)
![Build](https://img.shields.io/github/actions/workflow/status/cuhsat/fox/test.yaml?style=for-the-badge&label=build)
![Commits](https://img.shields.io/github/commit-activity/y/cuhsat/fox.svg?style=for-the-badge&label=commits)
![Release](https://img.shields.io/github/release/cuhsat/fox.svg?style=for-the-badge&label=release)

```console
go install github.com/cuhsat/fox/v4@latest
```

## Features
* Restricted Read-only access
* [Bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) detection
* Fast [Shannon entropy](https://en.wikipedia.org/wiki/Entropy_(information_theory)) calculation
* String carving and classification with
* Over 290 classes in [Hashcat](https://hashcat.net/wiki/doku.php?id=example_hashes) notation
* Dump Windows Shortcut and Prefetch files
* Dump [Linux ELF](https://refspecs.linuxfoundation.org/elf/elf.pdf) and [Windows PE/COFF](https://learn.microsoft.com/en-us/windows/win32/debug/pe-format) executables
* Check IPs, URLs and file hashes via the [VirusTotal API](https://www.virustotal.com/)
* Integral `grep`, `head`, `tail`, `hexdump`, `wc` like abilities
* Integral *Chain-of-Custody* receipt generation
* Many popular archive and compression formats
* Many popular cryptographic, fuzzy, image and fast hashes
* Special Hunt mode
  * Built-in support for [EnCase EWF](https://www.loc.gov/preservation/digital/formats/fdd/fdd000408.shtml) and raw `dd` images
  * Built-in log carving of [Linux Journals](https://systemd.io/JOURNAL_FILE_FORMAT/) and [Windows Event Logs](https://learn.microsoft.com/en-us/windows/win32/eventlog/event-log-file-format)
  * Built-in Super Timeline in [Common Event Format](https://www.microfocus.com/documentation/arcsight/arcsight-smartconnectors-8.3/cef-implementation-standard/Content/CEF/Chapter%201%20What%20is%20CEF.htm)
  * Built-in translation of over 51600 Event IDs
  * Built-in warning of critical system events
  * Filter events with [Sigma Rules](https://sigmahq.io/) syntax
  * Stream in [Splunk HEC](https://help.splunk.com/en/splunk-enterprise/leverage-rest-apis/rest-api-reference/10.0/input-endpoints/input-endpoint-descriptions) and [Elastic ECS](https://www.elastic.co/docs/reference/ecs) format
  * Save as `JSON`, `JSON Lines` or `SQLite3`

## Supports
File Formats
> evtx, journal, json, jsonl, lnk, pf, ELF, PE/COFF

Image Formats
> EWF-E01, EWF-S01, raw

Archive Formats
> 7zip, ar, CAB, cpio, RAR, RPM, tar, xar, ZIP

Compression Formats
> Brotli, bzip2, gzip, Kanzi, lz4, lzip, lzma, LZW, LZX, MinLZ, S2, Snappy, xz, zlib, zstd

Cryptographic Hashes
> BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512, MD2, MD4, MD5, MD6, RIPEMD-160, SHAKE128, SHAKE256, SHA1, SHA224, SHA256, SHA512, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Performance Hashes
> FNV-1, FNV-1a, Murmur3, XXH32, XXH64, XXH3

Similarity Hashes
> ImpHash, SSDeep, TLSH

Windows Specific
> LM, NT, PE Checksum  

Image Specific
> aHash, dHash, pHash

Checksums
> Adler32, Fletcher-4, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

## Manual
NAME
```
fox - The Forensic Examiners Swiss Army Knife
```
SYNOPSIS
```
fox [mode] [flags...] <paths...>
```
DESCRIPTION
```
Provides many useful features to leverage the forensic examination process.
```
MODES
```
cat    prints file (default)

hex    prints file in hex format

info   prints file infos and entropy

test   prints file test results

text   prints file text contents

hash   prints file hashes and checksums

hunt   hunt suspicious activities
```
FILE LIMITS
```
-h, --head
       limit head of file by ...

-t, --tail
       limit tail of file by ...

-n, --lines=NUMBER
       number of lines

-c, --bytes=NUMBER
       number of bytes
```
FILE LOADER
```
-p, --pass=PASSWORD
       password for decryption (7Z, RAR, ZIP)

-f, --file=FILENAME
       read extra paths from file

-i, --input=TEXT
       read input instead of file
```
LINE OUTPUT
```
-o, --output=FILE
       write all output to receipted file
```
LINE FILTER
```
-e, --regexp=PATTERN
       filter for lines that match pattern

-C, --context=NUMBER
       number of lines surrounding context of match

-B, --before=NUMBER
       number of lines leading context before match

-A, --after=NUMBER
       number of lines trailing context after match
```
PROFILE
```
-P, --profile=CORES
       parallel profile to use overall
```
DISABLE
```
-r, --raw
       don't process files at all

-q, --quiet
       don't print anything

    --no-file
       don't print filenames

    --no-line
       don't print line numbers

    --no-color
       don't colorize the output

    --no-pretty
       don't prettify the output

    --no-deflate
       don't deflate automatically

    --no-extract
       don't extract automatically

    --no-convert
       don't convert automatically

    --no-receipt
       don't write the receipt

    --no-warnings
       don't show any warnings
```
STANDARD
```
-d, --dry-run
       prints only the found filenames

-v, --verbose[=LEVEL]
       prints more details (v/vv/vvv)

    --version
       prints the version number
    
    --help
       prints this help message
```
POSITIONAL ARGUMENTS
```
Globbing paths to open or '-' to also read from STDIN.
```
EXAMPLES
```
$ fox -eWinlogon ./**/*.evtx
       Find occurrences in event logs

$ fox hex -hc512 disk.bin
       Show MBR in canonical hex

$ fox info -m0.9 ./**/*
       List high entropy files

$ fox test sample.exe
       Test suspicious file

$ fox text -w sample.exe
       Show strings in binary

$ fox hash -uTLSH files.7z
       Hash archive contents

$ fox hunt -sv ./**/*.E01
       Hunt down suspicious events
```
BUGS
```
Please submit any issues with fox to the project's bug tracker:
https://github.com/cuhsat/fox/issues
```
WWW
```
https://foxhunt.wtf
```
SEE ALSO
```
cat(1), grep(1), head(1), tail(1), hexdump(1), strings(1), wc(1)
```
---

*Disclaimer: This code was developed without the use of AI tooling and therefor does not contain any AI generated code, test or documentation. Furthermore, this code does not contain, employ or utilize AI tools in any other form. All data processed will not be shared with third parties except otherwise explicitly stated and permitted by the user.*

---
🦊 is released under the [GPL-3.0](LICENSE.md)