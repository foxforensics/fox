% FOX TEXT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **text** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints text contents.

FLAGS
=====

**-n, --min**=_length_

:   Minimum string _length_ (default: **3**).

**-x, --max**=_length_

:   Maximal string _length_ (default: **256**).

**-a, --ascii**

:   Shows only strings with ASCII encoding.

**-s, --sort**

:   Sorts strings alphabetically.

CLASSES
=======

**-w, --wtf**[=_level_]

:   Shows string classifications (w/ww/www).

**-F, --find**=_class_,...

:   Shows only strings that match _class_(es).

**-1, --first**

:   Shows only strings first class.

**-l, --list**

:   Shows only classification list.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox text -w sample.exe

:   Shows strings in binary.

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

**fox(1)**, **strings(1)**
