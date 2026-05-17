<div align="center">
  <img src="assets/img/fox.svg" width="160" alt=""/>
  <br/><br/>The Forensic Examiners Swiss Army Knife<br/><br/>

  [![Report](https://goreportcard.com/badge/github.com/f0x4n6/fox/v4?style=for-the-badge)](https://goreportcard.com/report/github.com/f0x4n6/fox/v4)
  [![Build](https://img.shields.io/github/actions/workflow/status/f0x4n6/fox/tests.yaml?style=for-the-badge&label=build)](https://github.com/f0x4n6/fox/actions)
  [![Release](https://img.shields.io/github/release/f0x4n6/fox.svg?style=for-the-badge&label=release)](https://github.com/f0x4n6/fox/releases)
</div>

## Abstract
Fox is a versatile commandline tool, built to support the examination process of file-based forensic evidence. It provides a wide spectrum of forensic capabilities in a cross-platform standalone binary.

## Features
* [x] Restricted read-only access
* [x] [Bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) detection
* [x] String carving and automatic classification
* [x] With 290+ classes in [Hashcat](https://hashcat.net/wiki/doku.php?id=example_hashes) notation
* [x] Parse Fortinet binary firewall logs
* [x] Parse Active Directory and other [EDB](https://learn.microsoft.com/en-us/windows/win32/extensible-storage-engine/extensible-storage-engine) files
* [x] Parse Windows shortcut and prefetch files
* [x] Parse [Linux ELF](https://refspecs.linuxfoundation.org/elf/elf.pdf) and [Windows PE/COFF](https://learn.microsoft.com/en-us/windows/win32/debug/pe-format) executables
* [x] Extract [Active Directory](https://learn.microsoft.com/en-us/troubleshoot/windows-server/windows-security/ntlm-user-authentication) hashes, users and computers  
* [x] Lookup NTLM hashes using 210000+ common passwords
* [x] Lookup URLs, IPs, domains and files via the [VirusTotal API](https://www.virustotal.com/)
* [x] Integral `grep`, `head`, `tail`, `uniq`, `wc`, `hexdump` like abilities
* [x] Integral syntax highlighting for many different formats
* [x] Integral fast [Shannon entropy](https://en.wikipedia.org/wiki/Entropy_(information_theory)) calculation
* [x] Integral *Chain-of-Custody* receipt generation
* [x] Support of path globbing and file streams
* [x] Support of encrypted `7z`, `Rar`, `Zip` archives
* [x] Many popular archive and compression formats
* [x] Many popular cryptographic, image, fuzzy and fast hashes
* [x] With [man pages](assets/man) for every command
* [x] Special [Hunt](assets/man/fox-hunt.md) command
  * [x] Built-in log carving of [Linux Journals](https://systemd.io/JOURNAL_FILE_FORMAT/) and [Windows Event Logs](https://learn.microsoft.com/en-us/windows/win32/eventlog/event-log-file-format)
  * [x] Built-in super timeline in [Common Event Format](https://www.microfocus.com/documentation/arcsight/arcsight-smartconnectors-8.3/cef-implementation-standard/Content/CEF/Chapter%201%20What%20is%20CEF.htm)
  * [x] Built-in translation of 51600+ event ids
  * [x] Built-in warning of critical system events
  * [x] Filter events with [Sigma Rules](https://sigmahq.io/) syntax
  * [x] Filter anomalies using [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
  * [x] Stream in [Splunk HEC](https://help.splunk.com/en/splunk-enterprise/leverage-rest-apis/rest-api-reference/10.0/input-endpoints/input-endpoint-descriptions) and [Elastic ECS](https://www.elastic.co/docs/reference/ecs) format
  * [x] Save as `JSON`, `JSON Lines`, `Parquet` or `SQLite` 

## Install
Install directly via the `go install` command:

```console
go install go.foxforensics.dev/fox/v4@latest
```

Standalone binaries and packages are available for:

|   OS    | Binaries                                                                                                                                                                       | Packages                                                                                                                                                                                                                                                                                                                                             |
|:-------:|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|  Linux  | [amd](https://foxforensics.dev/fox/releases/latest/download/fox_linux_amd64.tar.gz) \| [arm](https://foxforensics.dev/fox/releases/latest/download/fox_linux_arm64.tar.gz)     | [apk](https://foxforensics.dev/fox/releases/latest/download/fox_linux_amd64.apk) \| [deb](https://foxforensics.dev/fox/releases/latest/download/fox_linux_amd64.deb) \| [pkg](https://foxforensics.dev/fox/releases/latest/download/fox_linux_amd64.pkg.tar.zst) \| [rpm](https://foxforensics.dev/fox/releases/latest/download/fox_linux_amd64.rpm) |
|  macOs  | [amd](https://foxforensics.dev/fox/releases/latest/download/fox_darwin_amd64.tar.gz) \| [arm](https://foxforensics.dev/fox/releases/latest/download/fox_darwin_arm64.tar.gz)   | `brew install f0x4n6/fox/fox`                                                                                                                                                                                                                                                                                                                        |
| Windows | [amd](https://foxforensics.dev/fox/releases/latest/download/fox_windows_amd64.zip) \| [arm](https://foxforensics.dev/fox/releases/latest/download/fox_windows_arm64.zip)       | Binaries are UPX compressed                                                                                                                                                                                                                                                                                                                          |

## Examples

Find occurrences in event logs:
```console
fox -eWinlogon ./**/*.evtx
```

Show MBR in canonical hex:
```console
fox -hc512 disk.dd
```

Show NTLM password hashes:
```console
fox ad -LH NTDS.dit SYSTEM
```

Show all strings in a binary:
```console
fox str -w sample.exe
```

List only high entropy files:
```console
fox info -n6.0 ./**/*
```

Hash archive contents as MD5:
```console
fox hash -Amd5 files.7z
```

Hunt down critical events:
```console
fox hunt -u *.dd
```

## Capabilities
AD Records
> NTLM, Users, Computers

Log Formats
> EVTX, Journal, Fortigate

Binary Formats
> PE / COFF, ELF, ESE / EDB, LNK, PF

Archive Formats
> 7-Zip, AR, CAB, CFB, CPIO, ISO, MSI, RAR, RPM, TAR, XAR, ZIP

Compression Formats
> BGZF, Brotli, Bzip2, Gzip, Kanzi, LZ4, Lzip, LZMA, LZFSE, LZNT1, LZO, LZVN, LZW, LZX, MinLZ, S2, Snappy, XZ, zlib, zstd

Cryptographic Hashes
> BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512, GOST2012-256, GOST2012-512, HAS-160, LSH-256, LSH-512, MD2, MD4, MD5, MD6, RIPEMD-160, SHAKE128, SHAKE256, SHA1, SHA224, SHA256, SHA512, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512, Skein-224, Skein-256, Skein-384, Skein-512, SM3, Whirlpool

Performance Hashes
> DJB2, FNV-1, FNV-1a, Murmur3, RapidHash, SipHash, XXH32, XXH64, XXH3

Perceptual Hashes
> Average, Difference, Median, PHash, WHash, MarrHildreth, BlockMean, PDQ, RASH

Similarity Hashes
> ImpFuzzy, ImpHash, ImpHash0, SSDeep, TLSH

Windows Hashes
> LM, NT, PE

Checksums
> Adler32, Fletcher4, CRC16-CCITT, CRC32-C, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

Wordlists
> Default Passwords, Top Passwords, Worst Passwords, Corporate Passwords, Password Permutations, Common SSH Passwords, 100k Common English Passwords, 100k Common German Passwords, Medical Devices

---
🦊 is released under the [GPL-3.0](LICENSE.md)