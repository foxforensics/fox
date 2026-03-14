% FOX HEX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hex** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show file contents in hex format. This mode enforces the **--no-convert** flag.

FLAGS
=====

**-C, --canonical**

:   Format output as **canonical**.

**-H, --hexdump**

:   Format output like **hexdump**.

**-X, --xxd**

:   Format output like **xxd**.

Format Flags
------------

**-D, --decimal**

:   Format addresses as decimal.

Disable Flags
-------------

**-R, --no-format**

:   Don't format output at all.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox hex -hc512 disk.dd

:   Show MBR in canonical hex.

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

**fox(1)**, **hexdump(1)**, **xxd(1)**
