% FOX INFO(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **info** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show file infos with verdict. If the **--sort** flag is used, the files will be processed single-threaded, not _parallel_.

FLAGS
=====

**-s, --sort**

:   Sort files by path (slower).

**-b, --block**=_size_

:   Block _size_ for analysis.

Filter Flags
------------

**-n, --min**=_value_

:   Filter minimum entropy _value_ (default: **0.0**).

**-x, --max**=_value_

:   Filter maximal entropy _value_ (default: **1.0**).

Format Flags
------------

**-H, --human**

:   Format size in human-readable units.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ENVIRONMENT
===========

**FOX_API_KEY**

:   API key used to check files with the **VirusTotal** API. File checks will be conducted only by the files **SHA256** hash value, no files will be uploaded to VirusTotal (_optional_).

EXAMPLES
========

fox info -n0.8 ./**/*

:   List only high entropy files.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://github.com/cuhsat/fox/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxhunt.dev_>

SEE ALSO
========

**fox(1)**, **wc(1)**
