% FOX FIND(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **find** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints file infos and entropy.

FLAGS
=====

**-b, --block**=_size_

:   Block _size_ for calculations.

**-n, --min**=_value_

:   Minimum entropy _value_ (default: **0.0**).

**-x, --max**=_value_

:   Maximal entropy _value_ (default: **1.0**).

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox find -n0.9 ./**/*

:   Finds high entropy files

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://github.com/cuhsat/fox/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxhunt.wtf_>

SEE ALSO
========

**fox(1)**, **wc(1)**
