% FOX INFO(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **info** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show file infos and entropy. If the **--sort** flag is used, the files will be processed single-threaded, not _parallel_. If the **FOX_API_KEY** environment variable is set, then file hashes will be checked via the **VirusTotal** API. An extensive file report is available via the **--json** and **--jsonl** flags.

FLAGS
=====

**-s, --sort**

:   Sort files by path (slower).

**-j, --json**

:   Show infos as JSON objects.

**-J, --jsonl**

:   Show infos as JSON lines.

Block Flags
-----------

**-b, --block**=_size_

:   Block _size_ for analysis (default: all).

Filter Flags
------------

**-n, --min**=_value_

:   Filter minimum entropy _value_ (default: **0.0**).

**-x, --max**=_value_

:   Filter maximal entropy _value_ (default: **8.0**).

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ENVIRONMENT
===========

**FOX_API_KEY**

:   API key used to check files with the **VirusTotal** API. File checks will be conducted only by the files **SHA256** hash value, no files will be uploaded to VirusTotal. This environment variable is _optional_.

EXAMPLES
========

$ fox info -n6.0 ./**/*

:   List only high entropy files.

$ fox info -b1m db.sqlite3

:   List blocks by one MB size.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://foxforensics.dev/fox/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxforensics.dev/fox_>

SEE ALSO
========

**fox(1)**, **wc(1)**
