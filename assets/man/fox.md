% FOX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** \[_mode_] \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Fox provides many useful features to leverage your forensic examination process.

MODES
=====

cat

:   Prints contents (default mode).

hex

:   Prints contents in hex format.

info

:   Prints infos and entropy.

text

:   Prints text contents.

test

:   Prints test results.

hash

:   Prints hashes and checksums.

hunt

:   Hunts suspicious activities.

FLAGS
=====

**-i, --in**=_file_

:   Reads paths from _file_.

**-o, --out**=_file_

:   Writes output to _file_ (receipted).

Limit Flags
-----------

**-h, --head**

:   Limits head of file by **bytes** or **lines**.

**-t, --tail**

:   Limits tail of file by **bytes** or **lines**.

**-c, --bytes**=_number_

:   _Number_ of bytes.

**-l, --lines**=_number_

:   _Number_ of lines.

Filter Flags
------------

**-e, --regexp**=_pattern_

:   Filters output by _pattern_.

Crypto Flags
------------

**-P, --password**=_password_

:   Uses archive _password_ (only for **7Z**, **RAR**, **ZIP**).

Profile Flags
-------------

**-T, --threads**=_cores_

:   Uses parallel threads.

Disable Flags
-------------

**-r, --raw**

:   Don't process files at all.

**-q, --quiet**

:   Don't print anything.

**--no-file**

:   Don't print filenames.

**--no-line**

:   Don't print line numbers.

**--no-color**

:   Don't colorize the output.

**--no-pretty**

:   Don't prettify the output.

**--no-deflate**

:   Don't deflate automatically.

**--no-extract**

:   Don't extract automatically.

**--no-convert**

:   Don't convert automatically.

**--no-receipt**

:   Don't write the receipt.

**--no-warnings**

:   Don't show any warnings.

Standard Flags
--------------

**-p, --pause**

:   Prints only one page at a time.

**-d, --dry-run**

:   Prints only the found filenames.

**-v, --verbose**=[_level_]

:   Prints more details (v/vv/vvv).

**--version**

:   Prints the version number.

**--help**

:   Prints this help message.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ENVIRONMENT
===========

**FOX_VAR_NAME**

:   Global flags can be set through environment variables. The variable name must be prefixed with _FOX_ followed by an underscore and the flags name. All dots must be replaced with underscores.

EXAMPLES
========

fox -eWinlogon ./**/*.evtx

:   Finds occurrences in event logs.

fox hex -hc512 disk.bin

:   Shows MBR in canonical hex.

fox info -n0.9 ./**/*

:   Lists high entropy files

fox text -w sample.exe

:   Shows strings in binary.

fox test sample.exe

:   Tests suspicious file.

fox hash -uTLSH files.7z

:   Hashes archive contents.

fox hunt -sv ./**/*.E01

:   Hunts down suspicious events.

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

**cat(1)**, **grep(1)**, **head(1)**, **tail(1)**, **hexdump(1)**, **strings(1)**, **wc(1)**
