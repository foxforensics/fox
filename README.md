![fox](assets/fox_logo.png "fox")

The Forensic Swiss Army Knife. Providing many useful features to leverage your forensic examination process. Standalone binaries available for Windows, Linux and macOS.

![Status](https://img.shields.io/github/actions/workflow/status/cuhsat/fox/ci.yaml?style=flat-square&label=Status)
![Commits](https://img.shields.io/github/commit-activity/y/cuhsat/fox.svg?style=flat-square&label=Commits)
![Release](https://img.shields.io/github/release/cuhsat/fox.svg?style=flat-square&label=Release)

```console
go install github.com/cuhsat/fox/v4@latest
```

## Features
* Read-only filesystem access
* [Bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) detection
* Fast [Shannon entropy](https://en.wikipedia.org/wiki/Entropy_(information_theory)) calculation
* String carving and classification
* Integral `grep`, `head`, `tail`, `hexdump`, `wc` like abilities
* Hunt mode
  * Built-in file carving of [Linux Journals](https://systemd.io/JOURNAL_FILE_FORMAT/) and [Windows Event Logs](https://learn.microsoft.com/en-us/windows/win32/eventlog/event-log-file-format)
  * Built-in super timeline in [Common Event Format](https://www.microfocus.com/documentation/arcsight/arcsight-smartconnectors-8.3/cef-implementation-standard/Content/CEF/Chapter%201%20What%20is%20CEF.htm)
  * Built-in translation list of over 1500 Event IDs
  * Built-in warning of critical system events
  * Save as `JSON`, `JSON Lines` or `SQLite3`
* Supports
  * Over 160 string classes in [Hashcat](https://hashcat.net/wiki/doku.php?id=example_hashes) notation
  * Many popular archive and compression formats
  * Many popular cryptographic, fuzzy and fast hashes 
  * Data streaming in [Splunk HEC](https://help.splunk.com/en/splunk-enterprise/leverage-rest-apis/rest-api-reference/10.0/input-endpoints/input-endpoint-descriptions) or [Elastic ECS](https://www.elastic.co/docs/reference/ecs) format

## Usage
Type `fox --help` for more help:
```console
$ fox [COMMAND] [FLAGS] <PATHS>
```

## Examples
Find occurrences in event logs:
```console
$ fox cat -eWinlogon ./**/*.evtx
```

Show the MBR in canonical hex:
```console
$ fox hex -mc -hc512 disk.bin
```

Find ASCII strings in binaries:
```console
$ fox text -rwa8 sample.exe
```

List files with high entropy:
```console
$ fox info -a0.9 ./**/*
```

Hash the archive contents:
```console
$ fox hash -amd5,sha1 files.7z
```

Hunt down suspicious events:
```console
$ fox hunt -sxv ./**/*.dd
```

## Supports

File formats:
> evtx, journal, json, jsonl

Archive formats:
> 7zip, AR, CAB, CPIO, RAR, TAR, ZIP

Compression formats:
> Brotli, bzip2, gzip, Kanzi, lz4, lzip, lzma, LZW, LZX, MinLZ, S2, Snappy, xz, zlib, zstd

Cryptographic hashes:
> BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512
GOST-256, GOST-512, MD2, MD4, MD5, RIPEMD-160, SHA1, SHA256, SHA3
SHA3-224, SHA3-256, SHA3-384, SHA3-512, Tiger, Tiger2, Whirlpool

Performance hashes:
> FNV-1, FNV-1a, SipHash-64, SipHash-128, XXH64, XXH3

Similarity hashes:
> sdhash, SSDeep, TLSH

Checksums:
> Adler-32, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

## License
🦊 is released under the [GPL-3.0](LICENSE.md)