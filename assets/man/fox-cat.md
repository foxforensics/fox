% FOX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **cat** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints contents.

FLAGS
=====

**-C, --context**=_number_

:   Lines surrounding context of a match.

**-B, --before**=_number_

:   Lines leading context before a match.

**-A, --after**=_number_

:   Lines trailing context after a match.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox -eWinlogon ./**/*.evtx

:   Find occurrences in event logs.

SEE ALSO
========

**fox(1)**
