% FOX CHECK(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **check** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Check suspicious files using the VirusTotal API. An API key is required for this. No files will be uploaded to VirusTotal, file checks will be conducted only the files **SHA256** hash value.

FLAGS
=====

**-j, --json**

:   Show results as JSON objects.

**-J, --jsonl**

:   Show results as JSON lines.

Required
--------

**--api-key**=_apikey_

:   _API key_ for **VirusTotal** (hexadecimal format).

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ENVIRONMENT
===========

**FOX_API_KEY**

:   The **VirusTotal** API key can also be set through this environment variable.

EXAMPLES
========

fox check sample.exe

:   Check a suspicious file by hash.

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

**fox(1)**
