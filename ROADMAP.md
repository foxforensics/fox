# Roadmap

## 1. Bug Fixes
* ~~Done~~

## 2. Refactorings
* ~~Done~~

## 3. Features
* ~~Done~~

## 4. Optimizations
* ~~Done~~

## 5. Ideas
* Add advanced examples?
  > Mount the network share, decrypt the archive, extract the file into memory, calculate the hash and test it against known malware on VirusTotal:
  > ```console
  > fox test -pinfected \\lab\case\files.zip:ioc.exe
  > ```
* Add support for NTFS to disks in hunt mode
* Add support for MSI files
  * https://pkg.go.dev/github.com/asalih/go-msi
* Add animated console image to README
  * https://docs.asciinema.org/manual/cli/
* Add CEF syntax highlighting (pygments)
* Add global flag to change chroma style
* Add MSI release
  * https://dev.to/abdfnx/how-to-create-a-msi-file-for-go-program-je
* Add different endpoints for test
  * https://github.com/woanware/lookuper
  * HaveIBeenHacked for mails
* Add heatmap to hex mode
  * Custom BgColor from red to yellow in 256 steps?
* Add pager for output like moor?
  * https://github.com/walles/moor
* Add support for salted/keyed hashes?
  * `-S, --salt` flag
* Add Telf hash support?
  * https://github.com/trendmicro/telfhash
* Add LZXPRESS to own ESE fork?
  * https://forensics.wiki/compression/#lzxpress
  * https://github.com/Velocidex/go-prefetch/blob/master/lzxpress.go
  * https://github.com/fox-it/dissect.util/tree/main/dissect/util/compression
* Add LZNT1 deflate (magic bytes?)
  * https://github.com/Velocidex/go-ntfs/blob/d467c5e7dca0/lznt1.go
* Add vanity URL?
  * https://sagikazarmark.hu/blog/vanity-import-paths-in-go/
  * https://gitea.com/techknowlogick/go-vanity-url
* Add macOS trace3 log parser?
* Add Shellbag parser?
* Add Jumplist parser?
* Add specs to manpages?
  * `Fox Unified Event Format`
  * `Fox Unified Event Storage`
  * `Fox Chain of Custody Receipt`
