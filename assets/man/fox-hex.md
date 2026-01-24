% FOX HEX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hex** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints file contents in hex format.

FLAGS
=====

**-H, --hexdump**

:   Formats output like **hexdump**.

**-X, --xxd**

:   Formats output like **xxd**.

**-R, --raw**

:   Don't format the output.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox hex -hc512 disk.bin

:   Shows MBR in canonical hex.

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

**fox(1)**, **hexdump(1)**, **xxd(1)**
