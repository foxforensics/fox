% FOX TIME(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **time** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show the super timeline of different files and formats.

FLAGS
=====

**-b, --body**

:   Show timeline as body file version 3. All timestamps are in **UTC** time.

**-j, --json**

:   Show timeline as JSON objects.

**-J, --jsonl**

:   Show timeline as JSON lines.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive::file notation.

EXAMPLES
========

$ fox time "$MFT"

:   Show timeline of MFT entries.

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

**fox(1)**