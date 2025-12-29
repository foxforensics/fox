![fox](assets/logo.png "logo")

The Forensic Swiss Army Knife. Providing many useful features to leverage your forensic examination process. Standalone binaries available for Windows, Linux and macOS.

![Report](https://goreportcard.com/badge/github.com/cuhsat/fox/v4?style=for-the-badge)
![Build](https://img.shields.io/github/actions/workflow/status/cuhsat/fox/build.yaml?style=for-the-badge&label=build)
![Commits](https://img.shields.io/github/commit-activity/y/cuhsat/fox.svg?style=for-the-badge&label=commits)
![Release](https://img.shields.io/github/release/cuhsat/fox.svg?style=for-the-badge&label=release)

```console
go install github.com/cuhsat/fox/v4@latest
```

## Features
* Guaranteed read-only access
* [Bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) detection
* Fast [Shannon entropy](https://en.wikipedia.org/wiki/Entropy_(information_theory)) calculation
* Dumping of [Windows PE/COFF](https://learn.microsoft.com/en-us/windows/win32/debug/pe-format) executables
* String carving and classification
* Integral `grep`, `head`, `tail`, `hexdump`, `wc` like abilities
* Automatic Chain-of-Custody receipt generation
* Hunt mode
  * Built-in file carving of [Linux Journals](https://systemd.io/JOURNAL_FILE_FORMAT/) and [Windows Event Logs](https://learn.microsoft.com/en-us/windows/win32/eventlog/event-log-file-format)
  * Built-in super timeline in [Common Event Format](https://www.microfocus.com/documentation/arcsight/arcsight-smartconnectors-8.3/cef-implementation-standard/Content/CEF/Chapter%201%20What%20is%20CEF.htm)
  * Built-in translation of over 51600 Event IDs
  * Built-in warning of critical system events
  * Stream in [Splunk HEC](https://help.splunk.com/en/splunk-enterprise/leverage-rest-apis/rest-api-reference/10.0/input-endpoints/input-endpoint-descriptions) or [ECS](https://www.elastic.co/docs/reference/ecs) format
  * Save as `JSON`, `JSON Lines` or `SQLite3`
* Supports
  * Over 290 string classes in [Hashcat](https://hashcat.net/wiki/doku.php?id=example_hashes) notation
  * Many popular archive and compression formats
  * Many popular cryptographic, fuzzy and fast hashes

## Usage
Type `fox --help` for more help:
```console
$ fox [MODE] [FLAGS ...] <PATHS ...>
```

## Examples
Find occurrences in event logs:
```console
$ fox -eWinlogon ./**/*.evtx
```

Show the MBR in canonical hex:
```console
$ fox hex -hc512 disk.bin
```

List files with high entropy:
```console
$ fox info -m0.9 ./**/*
```

Find ASCII strings in binaries:
```console
$ fox text -rw sample.exe
```

Hash the archive contents:
```console
$ fox hash -Tmd5,sha1 files.7z
```

Hunt down suspicious events:
```console
$ fox hunt -sxv ./**/*.dd
```

## Supports

File formats:
> evtx, journal, json, jsonl, PE/COFF (.dll, .exe, .sys, ...)

Archive formats:
> 7zip, ar, CAB, cpio, RAR, RPM, tar, xar, ZIP

Compression formats:
> Brotli, bzip2, gzip, Kanzi, lz4, lzip, lzma, LZW, LZX, MinLZ, S2, Snappy, xz, zlib, zstd

Cryptographic hashes:
> BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512, MD2, MD4, MD5, MD6, SHA1, SHA256, SHA512, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512

Performance hashes:
> FNV-1, FNV-1a, Murmur3, XXH64, XXH3

Similarity hashes:
> SSDeep, TLSH

Windows hashes:
> LM, NT, PE

Checksums:
> Adler32, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

## Disclaimer
This code was developed without the use of AI tooling and therefor does not contain any AI generated code or documentation. Furthermore, this code does not contain, employ or utilize AI tooling in any other form. All data processed will not be shared with third parties under any circumstances.

🦊 is released under the [GPL-3.0](LICENSE.md).