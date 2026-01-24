% FOX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **info** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints infos and entropy.

FLAGS
=====

**-b, --block**=_size_

:   Block _size_ for calculations.

**-n, --min**=_value_

:   Minimum entropy _value_ (default: **0.0**).

**-x, --max**=_value_

:   Maximal entropy _value_ (default: **1.0**).

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox info -n0.9 ./**/*

:   Lists high entropy files

SEE ALSO
========

**fox(1)**, **wc(1)**
