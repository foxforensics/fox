% FOX STR(1) Version 5 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **str** \[_flags_ ...] \[**list** | _paths_ ...]

DESCRIPTION
===========

Show file string contents by carving the file(s). Please be aware that, using the **--sort** flag will buffer all found strings in memory. For large sets of data this could be very slow and take a serious amount of memory.

FLAGS
=====

**-a, --ascii**

:   Show only strings with ASCII encoding.

**-s, --sort**

:   Sort strings alphabetically.

**-t, --trim**

:   Trim strings whitespaces.

**-N, --min**=_length_

:   Minimum string _length_ (default: **3**).

**-X, --max**=_length_

:   Maximal string _length_ (default: **256**).

Class Flags
-----------

**-w, --what**[=_level_]

:   Show string classifications (**w**/**ww**/**www**).

**-C, --class**=_name_,...

:   Show only strings classes that match _name_. Implies **--what** flag at level _3_.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. If **list** is specified as _path_, only the list of the built-in classifications will be shown. To refer to paths inside archives, use the archive!file notation.

EXAMPLES
========

$ fox str -atN8 sample.exe

:   Show only long ASCII strings.

$ fox str -wCurl sample.exe

:   Show all URLs in a binary.

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

**fox(1)**, **strings(1)**
