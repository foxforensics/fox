% FOX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hex** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints contents in hex format.

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

SEE ALSO
========

**fox(1)**, **hexdump(1)**, **xxd(1)**
