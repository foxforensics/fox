% FOX INFO(1) Version 5 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **info** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show file infos and entropy.

FLAGS
=====

**-j, --json**

:   Show infos as JSON objects.

**-J, --jsonl**

:   Show infos as JSON lines.

Block Flags
-----------

**-B, --block**=_size_

:   Block _size_ for analysis (default: all). The _size_ can be either defined as raw bytes or with a size suffix.

Filter Flags
------------

**-N, --min**=_value_

:   Filter minimum entropy _value_ (default: **0.0**).

**-X, --max**=_value_

:   Filter maximal entropy _value_ (default: **8.0**).

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive::file notation.

EXAMPLES
========

$ fox info -N6.0 ./

:   List only high entropy files.

$ fox info -B1m backup.mdf

:   List blocks by one megabyte.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://foxforensics.eu/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxforensics.eu_>

SEE ALSO
========

**fox(1)**, **wc(1)**, **sha256sum(1)**
