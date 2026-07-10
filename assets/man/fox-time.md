% FOX TIME(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **time** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show super timeline of MFT, LNK or PF files. Please be aware that, using the **--sort** flag will buffer all found events in memory. For large sets of data this could be very slow and take a serious amount of memory. All timestamps will be normalized to **UTC**.

FLAGS
=====

**-s, --sort**

:   Sort timeline chronologically.

**-j, --json**

:   Show timeline as JSON objects.

**-J, --jsonl**

:   Show timeline as JSON lines.

Format Flags
------------

**-b, --bodyfile**

:   Show in **Body File Version 3** format.

**-t, --timesketch**

:   Show in **Timesketch** format.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive::file notation.

EXAMPLES
========

$ fox time -b ./$MFT

:   Show MFT entries as body file.

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