% FOX TEXT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **text** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Show file text contents by carving the file(s). This mode enforces the **--no-convert** flag.

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

Format Flags
------------

**-D, --decimal**

:   Format addresses as decimal.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox text -w ioc.exe

:   Show strings in binary.

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

**fox(1)**, **strings(1)**
