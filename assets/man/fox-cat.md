% FOX CAT(1) Version 4 | Fox Documentation

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

:   _Number_ of lines surrounding context of a match.

**-B, --before**=_number_

:   _Number_ of lines leading context before a match.

**-A, --after**=_number_

:   _Number_ of lines trailing context after a match.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox -eWinlogon ./**/*.evtx

:   Finds occurrences in event logs.

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

**fox(1)**, **cat(1)**, **grep(1)**, **less(1)**
