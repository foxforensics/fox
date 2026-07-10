<div align="center">
  <img src="assets/img/fox.svg" width="160" alt=""/>
  <br/><br/><b>The Forensic Examiners Swiss Army Knife</b><br/><br/>

  [![Build](https://img.shields.io/github/actions/workflow/status/foxforensics/fox/tests.yaml?style=for-the-badge&label=build)](https://github.com/foxforensics/fox/actions)
  [![Release](https://img.shields.io/github/release/foxforensics/fox.svg?style=for-the-badge&label=release)](https://github.com/foxforensics/fox/releases)

  Fox is a powerful CLI tool, built to support the examination process of file-based forensic artifacts.</br>
  It provides a wide spectrum of forensic capabilities in a cross-platform standalone binary.
</div>

## Features
* [x] Restricted read-only access
* [x] [Bidirectional character](https://nvd.nist.gov/vuln/detail/CVE-2021-42574) detection
* [x] String carving and automatic classification
* [x] With 290+ classes in [Hashcat](https://hashcat.net/wiki/doku.php?id=example_hashes) notation
* [x] Parse Fortinet binary firewall logs
* [x] Parse MFT, LNK, PF, PST binary files
* [x] Parse Active Directory and other [EDB](https://learn.microsoft.com/en-us/windows/win32/extensible-storage-engine/extensible-storage-engine) files 
* [x] Parse [Linux ELF](https://refspecs.linuxfoundation.org/elf/elf.pdf) and [Windows PE/COFF](https://learn.microsoft.com/en-us/windows/win32/debug/pe-format) executables
* [x] Extract [Active Directory](https://learn.microsoft.com/en-us/troubleshoot/windows-server/windows-security/ntlm-user-authentication) hashes, users, groups, computers  
* [x] Lookup NTLM hashes using 210000+ entry wordlists
* [x] Integral `grep`, `head`, `tail`, `uniq`, `wc`, `hexdump` like abilities
* [x] Integral syntax highlighting for many different formats
* [x] Integral super timeline with timestomp checks
* [x] Integral fast [Shannon entropy](https://en.wikipedia.org/wiki/Entropy_(information_theory)) calculation
* [x] Integral *Chain-of-Custody* receipt generation
* [x] Support of path globbing and file streams
* [x] Support of archive and file path sanitization
* [x] Support of encrypted `7z`, `Rar`, `Zip` archives
* [x] Many popular archive and compression formats
* [x] Many popular cryptographic, image, fuzzy, fast and OS hashes
* [x] Bundled [man pages](assets/man) for every command
* [x] Advanced [Hunt](assets/man/fox-hunt.md) command
  * [x] Built-in log carving of [Linux Journals](https://systemd.io/JOURNAL_FILE_FORMAT/) and [Windows Event Logs](https://learn.microsoft.com/en-us/windows/win32/eventlog/event-log-file-format)
  * [x] Built-in super timeline in [Common Event Format](https://www.microfocus.com/documentation/arcsight/arcsight-smartconnectors-8.3/cef-implementation-standard/Content/CEF/Chapter%201%20What%20is%20CEF.htm)
  * [x] Built-in translation of 51600+ event ids
  * [x] Built-in warning of critical system events
  * [x] Filter events with [Sigma Rules](https://sigmahq.io/) syntax
  * [x] Stream events in [Splunk](https://help.splunk.com/en/splunk-enterprise/leverage-rest-apis/rest-api-reference/10.0/input-endpoints/input-endpoint-descriptions) or [Elastic](https://www.elastic.co/docs/reference/ecs) format
  * [x] Save as `JSON`, `JSON Lines` or `Parquet` 

## Install
Install the development version directly via `go`:

```console
go install go.foxforensics.eu/fox/v5@latest
```

Standalone binaries and packages are available for:

|   OS    | Binaries                                                                                                                                                                                   | Packages                                                                                                                                                                                                                                                                                                                                                                         |
|:-------:|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|  Linux  | [amd](https://github.com/foxforensics/fox/releases/latest/download/fox_linux_amd64.tar.gz) \| [arm](https://github.com/foxforensics/fox/releases/latest/download/fox_linux_arm64.tar.gz)   | [apk](https://github.com/foxforensics/fox/releases/latest/download/fox_linux_amd64.apk) \| [deb](https://github.com/foxforensics/fox/releases/latest/download/fox_linux_amd64.deb) \| [pkg](https://github.com/foxforensics/fox/releases/latest/download/fox_linux_amd64.pkg.tar.zst) \| [rpm](https://github.com/foxforensics/fox/releases/latest/download/fox_linux_amd64.rpm) |
|  macOs  | [amd](https://github.com/foxforensics/fox/releases/latest/download/fox_darwin_amd64.tar.gz) \| [arm](https://github.com/foxforensics/fox/releases/latest/download/fox_darwin_arm64.tar.gz) | `brew install foxforensics/fox/fox`                                                                                                                                                                                                                                                                                                                                              |
| Windows | [amd](https://github.com/foxforensics/fox/releases/latest/download/fox_windows_amd64.zip) \| [arm](https://github.com/foxforensics/fox/releases/latest/download/fox_windows_arm64.zip)     | Binaries are standalone executables                                                                                                                                                                                                                                                                                                                                              |

## Examples

Find occurrences in event logs:
```console
fox -FWinlogon ./**/*.evtx
```

Show MBR in canonical hex:
```console
fox -L512b image.dd
```

Show NTLM password hashes:
```console
fox ad -hl NTDS.dit SYSTEM
```

Show all strings in a binary:
```console
fox str -w sample.exe
```

List only high entropy files:
```console
fox info -N6.0 ./
```

Show entries as body file:
```console
fox time -b ./$MFT
```

Hash archive contents as MD5:
```console
fox hash -Hmd5 files.7z
```

Hunt down critical events:
```console
fox hunt -u *.dd
```

## Capabilities
AD Records
> NTLM, Users, Groups, Computers

Log Formats
> EVTX, Journal, Fortigate

Binary Formats
> PE / COFF, ELF, ESE / EDB, MFT, LNK, PF, PST

Archive Formats
> 7-Zip, AR, CAB, CFB, CPIO, ISO, MSI, RAR, RPM, TAR, XAR, ZIP

Compression Formats
> BGZF, Brotli, Bzip2, Gzip, Kanzi, LZ4, Lzip, LZMA, LZFSE, LZO, LZVN, LZW, LZX, MinLZ, S2, Snappy, XZ, zlib, zstd

Cryptographic Hashes
> BLAKE2S-256, BLAKE2B-256, BLAKE2B-384, BLAKE2B-512, BLAKE3-256, BLAKE3-512, GOST2012-256, GOST2012-512, HAS-160, LSH-256, LSH-512, MD2, MD4, MD5, MD6, RIPEMD-160, SHAKE128, SHAKE256, SHA1, SHA224, SHA256, SHA512, SHA3, SHA3-224, SHA3-256, SHA3-384, SHA3-512, Skein-224, Skein-256, Skein-384, Skein-512, SM3, Whirlpool

Performance Hashes
> DJB2, FNV-1, FNV-1a, Murmur3, RapidHash, SipHash, XXH32, XXH64, XXH3

Perceptual Hashes
> Average, Difference, Median, PHash, WHash, MarrHildreth, BlockMean, PDQ, RASH

Similarity Hashes
> ImpFuzzy, ImpHashO, ImpHashS, sdhash, SSDeep, TLSH

Windows Specific
> LM, NT, PE

Unix Specific
> BSD, ELF, SYSV

Checksums
> Adler32, Fletcher4, Luhn, CRC16-CCITT, CRC32-C, CRC32-K, CRC32-IEEE, CRC64-ECMA, CRC64-ISO

Wordlists
> English / German / French / Spanish Common Passwords, Most Common Passwords, Most Used Passwords, Default Passwords, Corporate Passwords, Production Passwords, Milw0rm Dictionary, Conficker Dictionary, Medical Devices, Seasons 

## Building
To build the lastest development version execute:
```console
go build && ./fox --version
```

---
🦊 is released under the [GPL-3.0](LICENSE.md). All code is entirely written by human authors.
