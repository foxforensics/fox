% FOX CAT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **cat** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints file contents.

FLAGS
=====

**-u, --uniq**

:   Filters by unique hash (**XXH3**).

**-D, --dist**=_length_

:   Filters by Levenshtein distance (slow).

**-e, --regexp**=_pattern_

:   Filters by regular expression.

RegExp Flags
------------

**-C, --context**=_lines_

:   _Lines_ surrounding a match.

**-B, --before**=_lines_

:   _Lines_ leading before a match.

**-A, --after**=_lines_

:   _Lines_ trailing after a match.

Syntax Flags
------------

**-S, --syntax**=_type_

:   Force syntax highlighting _type_.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox -eWinlogon ./**/*.evtx

:   Find occurrences in event logs.

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

**fox(1)**, **cat(1)**, **grep(1)**, **uniq(1)**
