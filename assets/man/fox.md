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

**s, str**

:   Show file contained strings.

**i, info**

:   Show file infos and entropy.

**h, hash**

:   Show file hashes and checksums.

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

Unique Flags
------------

**-u, --uniq**

:   Filter using unique hash (**XXH3**).

**-D, --dist**=_length_

:   Filter using Levenshtein distance.

Filter Flags
------------

**-e, --regexp**=_pattern_

:   Filter output by _pattern_.

**-C, --context**=_lines_

:   _Lines_ surrounding a match.

**-B, --before**=_lines_

:   _Lines_ leading before a match.

**-A, --after**=_lines_

:   _Lines_ trailing after a match.

Special Flags
-------------

**-p, --password**=_password_

:   Use archive _password_ (only for **7Z**, **RAR**, **ZIP**).

**-P, --parallel**=_cores_

:   Use _cores_ for parallel processing.

Display Flags
-------------

**-T, --force-text**

:   Force output as text (in default mode).

**-X, --force-hex**

:   Force output as hex (in default mode).

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

**FOX_LEXER**

:   Force syntax highlighting lexer. Only available in default mode.

**FOX_STYLE**

:   Force syntax highlighting style. Only available in default mode.

EXAMPLES
========

fox -eWinlogon ./**/*.evtx

:   Find occurrences in event logs.

fox -hc512 disk.dd

:   Show MBR in canonical hex.

fox str -w sample.exe

:   Show all strings in a binary.

fox info -n6.0 ./**/*

:   List only high entropy files.

fox hash -Amd5 files.7z

:   Hash archive contents as MD5.

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

**cat(1)**, **grep(1)**, **head(1)**, **tail(1)**, **uniq(1)**, **wc(1)**, **strings(1)**, **hexdump(1)**, **sha256sum(1)**
