% FOX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** \[_command_] \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Fox is a CLI tool, build to support the examination process of file based forensic artifacts, by providing the most useful features in a cross-platform standalone binary. As with any Swiss Army knife, there are many specific power tools that offer more in-depth functionality, but sometimes all you need is a simple screwdriver.

All files will only be processed read-only. A Chain-of-Custody receipt is generated upon every file output.

COMMANDS
========

If no command is passed, then `cat` will be used by default.

File Commands
-------------

**c, cat**

:   Show file contents (default).

**x, hex**

:   Show file contents in hex format.

**s, str**

:   Show file contained strings.

**l, stat**

:   Show file stats and entropy.

**h, hash**

:   Show file hashes and checksums.

Misc Commands
-------------

**v, check**

:   Check suspicious items online.

**d, dump**

:   Dump Active Directory secrets.

**e, hunt**

:   Hunt critical system events.

FLAGS
=====

**-i, --in**=_file_

:   Read paths from _file_.

**-o, --out**=_file_

:   Write output to _file_ (receipted).

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

Special Flags
-------------

**-p, --password**=_password_

:   Use archive _password_ (only for **7Z**, **RAR**, **ZIP**).

**-T, --threads**=_cores_

:   Use parallel threads.

Disable Flags
-------------

**-r, --raw**

:   Don't process files at all.

**-q, --quiet**

:   Don't print anything.

**-N, --no-pretty**

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

Standard Flags
--------------

**-v, --verbose**[=_level_]

:   Print more details (v/vv/vvv).

**-d, --dry-run**

:   Print only the found files.

**--version**

:   Print the version number.

**--help**

:   Print this help message.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive:file notation. 

ENVIRONMENT
===========

**FOX_VAR_NAME**

:   Global flags can be set through environment variables. The variable name must be prefixed with _FOX_ followed by an underscore and the flags name. All dots must be replaced with underscores. A general proxy server can be set through the environment variables _HTTPS_PROXY_, _HTTP_PROXY_ and _NO_PROXY_.

EXAMPLES
========

fox -eWinlogon ./**/*.evtx

:   Find occurrences in event logs.

fox hex -hc512 disk.dd

:   Show MBR in canonical hex.

fox str -w sample.exe

:   Show all strings in a binary.

fox stat -n0.8 ./**/*

:   List only high entropy files.

fox hash -Amd5 files.7z

:   Hash archive contents as MD5.

fox check sample.exe

:   Check a suspicious file by hash.

fox dump system ntds.dit

:   Dump users and password hashes.

fox hunt -u *.dd

:   Hunt down critical events.

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
