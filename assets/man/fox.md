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

:   Shows file contents (default mode).

hex

:   Shows file contents in hex format.

text

:   Shows file contained strings.

hash

:   Shows file hashes and checksums.

list

:   Lists file infos and entropy.

test

:   Tests suspicious files.

hunt

:   Hunts suspicious events.

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

Archive Flags
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

**--no-syntax**

:   Don't colorize the syntax.

**--no-pretty**

:   Don't prettify the output.

**--no-strict**

:   Don't stop on parser errors.

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

**-d, --dry-run**

:   Prints only the found files.

**-v, --verbose**[=_level_]

:   Prints more details (v/vv/vvv).

**--version**

:   Prints the version number.

**--help**

:   Prints this help message.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive:file notation. 

ENVIRONMENT
===========

**FOX_VAR_NAME**

:   Global flags can be set through environment variables. The variable name must be prefixed with _FOX_ followed by an underscore and the flags name. All dots must be replaced with underscores.

EXAMPLES
========

fox -eWinlogon ./**/*.evtx

:   Find occurrences in event logs.

fox hex -hc512 disk.dd

:   Show MBR in canonical hex.

fox text -w ioc.exe

:   Show strings in binary.

fox hash -Amd5 files.7z

:   Hash archive contents.

fox list -n0.9 ./**/*

:   List high entropy files.

fox test ioc.exe

:   Test a suspicious file.

fox hunt -sv ./**/*.dd

:   Hunt down suspicious events.

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

**cat(1)**, **grep(1)**, **head(1)**, **tail(1)**, **uniq(1)**, **wc(1)**, **hexdump(1)**, **strings(1)**
