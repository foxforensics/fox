% FOX TEST(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **test** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Test suspicious files using VirusTotal. An API key is required for this. This mode enforces the **--no-convert** flag. This mode will return **3** as exit code, if any test returns at least one alerting result.

FLAGS
=====

**-D, --domain**

:   File(s) contains _domains_.

**-U, --url**

:   File(s) contains _urls_.

**-I, --ip**

:   File(s) contains _ips_.

Required
--------

**-k, --key**=_apikey_

:   **VirusTotal** API _key_.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ENVIRONMENT
===========

**FOX_KEY**

:   The **VirusTotal** API key can also be set through this environment variable.

EXAMPLES
========

fox test ioc.exe

:   Test a suspicious file.

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
