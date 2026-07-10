% FOX(1) Version 5 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** \[_command_] \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Fox is a powerful CLI tool, built to support the examination process of file-based forensic artifacts. It provides a wide spectrum of forensic capabilities in a cross-platform standalone binary. All files will only be processed read-only. A Chain-of-Custody receipt is generated upon every file output.

COMMANDS
========

If no command is passed, then the **cat** command will be used by default and the file contents will be shown as either text or hex output according to their content type.

**a, ad**

:   Show Active Directory infos.

**c, cat**

:   Show file contents (default).

**s, str**

:   Show file contained strings.

**i, info**

:   Show file infos and entropy.

**t, time**

:   Show file super timeline.

**h, hash**

:   Show file hashes and checksums.

**x, hunt**

:   Hunt critical system events.

FLAGS
=====

**-I, --in**=_file_

:   Read paths from _file_. **~** expansion is applied.

**-O, --out**=_file_

:   Write output to _file_. A receipt of the written file will be created automatically alongside.

Filter Flags
------------

**-L, --limit**=_number_

:   Filter using byte or line count. The value can be either specified as decimal (_b_), hexadecimal (_h_) or line count (_l_). A positive value implies only to show the leading _number_ of bytes or lines. A negative value implies only to show the trailing _number_ of bytes or lines. 

**-F, --find**=_regex_

:   Filter using regular expression _regex_. Regular expressions do not have constant time guarantees and allow backtracking. All regular expressions are PCRE-compatible with .NET and Perl5.

Process Flags
-------------

**-T, --threads**=_cores_

:   Use _cores_ for parallel threads. The default is the number of logical CPUs available for the process. This flag only affects file processing, not file discovery.

**-P, --password**=_text_

:   Use _text_ as password to decrypt encrypted archives. Use '-' to read the password securely from **STDIN(4)**. The password can also be set as environment variable. This is only supported for **7z**, **Rar** and **Zip** archives.

Disable Flags
-------------

**-r, --raw**[=_level_]

:   Don't process files (**r**/**rr**/**rrr**/**rrrr**). Level _1_ implies **--no-pretty**. Level _2_ implies **--no-convert**. Level _3_ implies **--no-deflate** and **--no-extract** additional. Level _4_ implies **--no-protect** additional.

**-q, --quiet**

:   Don't print anything.

**--no-pretty**

:   Don't prettify the output (**-r**). Header line and colors will be turned off.

**--no-convert**

:   Don't convert automatically (**-rr**).

**--no-deflate**

:   Don't deflate automatically (**-rrr**).

**--no-extract**

:   Don't extract automatically (**-rrr**).

**--no-protect**

:   Don't apply any loader protection (**-rrrr**). This effectively removes all resource limiting safeguards from the loading pipeline, which are in place to mitigate malicious crafted files like zip bombs.

**--no-receipt**

:   Don't write the receipt.

Standard Flags
--------------

**-v, --verbose**[=_level_]

:   Print more details (**v**/**vv**).

**-d, --dry-run**

:   Print only the found files.

**--version**

:   Print the version number.

**--help**

:   Print this help message.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. To refer to paths inside archives, use the archive::file notation. 

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

$ fox -FWinlogon ./**/*.evtx

:   Show occurrences in event logs.

$ fox -L512b image.dd

:   Show MBR in canonical hex.

$ fox ad NTDS.dit SYSTEM

:   Show NTLM password hashes.

$ fox str -w sample.exe

:   Show all strings in a binary.

$ fox info -N6.0 ./

:   List only high entropy files.

$ fox time -b ./$MFT

:   Show entries as body file.

$ fox hash -Hmd5 files.7z

:   Hash archive contents as MD5.

$ fox hunt -u *.dd

:   Hunt down critical events.

$ fox help info

:   Show help on sub commands.

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

**cat(1)**, **grep(1)**, **head(1)**, **tail(1)**, **sort(1)**, **uniq(1)**, **wc(1)**, **strings(1)**, **hexdump(1)**, **sha256sum(1)**
