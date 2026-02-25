% FOX STAT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **stat** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Shows file stats and entropy. If the **--sort** flag is used, the files will be processed single-threaded. This mode enforces the **--no-convert** flag.

FLAGS
=====

**-s, --sort**

:   Sorts files by path (slower).

**-b, --block**=_size_

:   Uses block _size_ for analysis.

Filter Flags
------------

**-n, --min**=_value_

:   Filters for minimum entropy _value_ (default: **0.0**).

**-x, --max**=_value_

:   Filters for maximal entropy _value_ (default: **1.0**).

Format Flags
------------

**-h, --human**

:   Formats size in human-readable units.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox stat -n0.9 ./**/*

:   List high entropy files

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
