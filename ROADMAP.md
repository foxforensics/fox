# Roadmap

## 1. Bug Fixes
* ~~Done~~

## 2. Refactorings
* ~~Done~~

## 3. Features
* Add Sigma rules to hunt mode
  * -r, --rules=FILE, []*os.File
  * event.Match(rules)
  * Warning for incompatible rules
  * Built-in rules for Windows and Linux 
  * https://github.com/bradleyjkemp/sigma-go
  * https://github.com/goccy/go-yaml

## 4. Optimizations
* Add host, user, time to -f *.coc
* Add Unicode to string carving
  * https://people.cs.umass.edu/~liberato/courses/2019-spring-compsci365/lecture-notes/04-carving-strings-and-unicode/
* SMap speed
  * https://dev.to/moseeh_52/efficient-file-reading-in-go-mastering-bufionewscanner-vs-osreadfile-4h05
  * https://dave.cheney.net/high-performance-json.html

## 5. Ideas
* Add Heatmap to hex command?
* Use reflow algos?
  * https://github.com/muesli/reflow
