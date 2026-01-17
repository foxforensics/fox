# Roadmap

## 1. Bug Fixes
* ~~Done~~

## 2. Refactorings
* ~~Done~~

## 3. Features
* ~~Done~~

## 4. Optimizations
* Optimize for files larger than memory
  * All modes use heaps sequentially
  * Add stream instead of head
    * Use io.Reader internally
    * Peek for magic bytes
      * use io.TeeReader? or io.ReadSeeker?
      * magic.Detect -> io.Reader?
  * Move SMap from Heaps to Text/Hex buffer
    * What about -n lines parameter?
    * Solve this with counting of line breaks?
  * What about archives?

## 5. Ideas
* Add Fortigate Firewall log parser?
* Add macOS trace3 log parser?
* Add Shellbag parser?
* Add Jumplist parser?
* Add specs?
  * `Fox Unified Event Format`
  * `Fox Unified Event Storage`
  * `Fox Chain of Custody Receipt`
