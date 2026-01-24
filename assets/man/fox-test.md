% FOX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **test** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints test results.

FLAGS
=====

**-k, --key**=_apikey_

:   Set key for **VirusTotal** API.

**-U, --url**=_url_,...

:   Test suspicious URL.

**-I, --ip**=_ip_,...

:   Test suspicious IP.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ENVIRONMENT
===========

**FOX_KEY**

:   The **VirusTotal** API key can also be set through this environment variable.

EXAMPLES
========

fox test sample.exe

:   Test suspicious file.

SEE ALSO
========

**fox(1)**
