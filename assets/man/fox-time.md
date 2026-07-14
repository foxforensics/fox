% FOX TIME(1) Version 5 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **time** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show the super timeline of MFT, LNK or PF files as bodyfile version 3. Please be aware that, using the **--sort** flag will buffer all found events in memory. For large sets of data this could be very slow and take a serious amount of memory. All timestamps will be normalized to **UTC**.

FLAGS
=====

**-s, --sort**

:   Sort timeline chronologically.

**-j, --json**

:   Show timeline as JSON objects.

**-J, --jsonl**

:   Show timeline as JSON lines.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive!file notation.

EXAMPLES
========

$ fox time ./$MFT

:   Show MFT entries as bodyfile.

$ fox time -s ./**/*.pf

:   Show entries chronologically.

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

**fox(1)**, **sort(1)**