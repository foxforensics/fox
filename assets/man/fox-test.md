% FOX TEST(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **test** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Prints file test results. A VirusTotal API key is required.

FLAGS
=====

**-k, --key**=_apikey_

:   Sets **VirusTotal** API _key_.

**-U, --url**=_url_,...

:   Tests suspicious _url_.

**-I, --ip**=_ip_,...

:   Tests suspicious _ip_.

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

:   Tests suspicious file.

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

**fox(1)**
