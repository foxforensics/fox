% FOX CAT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **cat** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show content of file(s) as text or hexdump. This command will be used as default, if no other command is specified.

FLAGS
=====

**-t, --text**

:   Force output exclusively as text (in default mode).

**-x, --hex**

:   Force output exclusively as hex (in default mode).

Unique Flags
------------

**-u, --uniq**

:   Unique by **XXH3** hash sum.

**-D, --dist**=_length_

:   Unique by Levenshtein distance.

Filter Flags
------------

**-F, --find**=_pattern_

:   Filter using regular expression _pattern_. Regular expressions do not have constant time guarantees and allow backtracking. All regular expressions are compatible with Perl5 and .NET.

**-C, --context**=_lines_

:   _Lines_ surrounding a match. Includes **-B** and **-A** flag.

**-B, --before**=_lines_

:   _Lines_ leading before a match.

**-A, --after**=_lines_

:   _Lines_ trailing after a match.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive:file notation.

EXAMPLES
========

$ fox cat -FWinlogon ./**/*.evtx

:   Show occurrences in event logs.

$ fox cat -L512b image.dd

:   Show MBR in canonical hex.

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

**cat(1)**, **grep(1)**, **uniq(1)**, **hexdump(1)**