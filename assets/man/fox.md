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

:   Hunt suspicious activities.

FLAGS
=====

File Flags
----------

**-i, --in**=_file_

:   Read paths from file.

**-o, --out**=_file_

:   Write output to file (receipted).

Limit Flags
-----------

**-h, --head**

:   Limit head of file by **bytes** or **lines**.

**-t, --tail**

:   Limit tail of file by **bytes** or **lines**.

**-c, --bytes**=_number_

:   _Number_ of bytes.

**-l, --lines**=_number_

:   _Number_ of lines.

Filter Flags
------------

**-e, --regexp**=_pattern_

:   Filter output by _pattern_.

Crypto Flags
------------

**-p, --password**=_password_

:   Archive _password_ (only for _7Z_, _RAR_, _ZIP_).

Profile Flags
-------------

**-P, --parallel**=_cpus_

:   Parallel processing usage.

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

**-m, --pause**

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

:   Find occurrences in event logs.

fox hex -hc512 disk.bin

:   Show MBR in canonical hex.

fox info -n0.9 ./**/*

:   List high entropy files

fox text -w sample.exe

:   Show strings in binary.

fox test sample.exe

:   Test suspicious file.

fox hash -uTLSH files.7z

:   Hash archive contents.

fox hunt -sv ./**/*.E01

:   Hunt down suspicious events.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://github.com/cuhsat/fox/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxhunt.wtf_>

AUTHOR
======

Christian Uhsat <fox at foxhunt dot wtf>

SEE ALSO
========

**cat(1)**, **grep(1)**, **head(1)**, **tail(1)**, **more(1)**, **hexdump(1)**, **strings(1)**, **wc(1)**
