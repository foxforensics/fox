% FOX(1) Version 4 | Fox Documentation

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

**-n, --min**=_number_

:   Minimum string length (default: 3).

**-x, --max**=_number_

:   Maximal string length (default: 256).

**-a, --ascii**

:   Show only strings with ASCII encoding.

**-s, --sort**

:   Sort strings alphabetically.

CLASSES
=======

**-w, --wtf**[=_level_]

:   Show string classifications (w/ww/www).

**-F, --find**=_class_,...

:   Show only strings that match class(es).

**-1, --first**

:   Show only strings first class.

**-l, --list**

:   Show only classification list.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox text -w sample.exe

:   Show strings in binary.

SEE ALSO
========

**fox(1)**
