% FOX STR(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **str** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show file string contents by carving the file(s).

FLAGS
=====

**-n, --min**=_length_

:   Minimum string _length_ (default: **3**).

**-x, --max**=_length_

:   Maximal string _length_ (default: **256**).

**-a, --ascii**

:   Show only strings with ASCII encoding.

**-s, --sort**

:   Sort strings alphabetically.

Class Flags
-----------

**-w, --wtf**[=_level_]

:   Show string classifications (w/ww/www).

**-F, --find**=_class_,...

:   Show only strings that match _class_(es).

**-1, --first**

:   Show only strings first class.

**-L, --list**

:   Show only classification list.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox str -w sample.exe

:   Show all strings in a binary.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://foxforensics.dev/fox/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxforensics.dev/fox_>

SEE ALSO
========

**fox(1)**, **strings(1)**
