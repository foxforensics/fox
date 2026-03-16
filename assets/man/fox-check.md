% FOX CHECK(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **check** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Check suspicious files, domains, mails, URLs and IPs using the HaveIBeenPwned or VirusTotal API. An API key is required for this. This command enforces the **--no-convert** flag. No files will be uploaded to VirusTotal, file checks will be conducted only the files **SHA256** hash value.

FLAGS
=====

**-j, --json**

:   Show results as JSON objects.

**-J, --jsonl**

:   Show results as JSON lines.

Content Flags
-------------

**-D, --domain**

:   File(s) contains a list of _domains_ separated by line breaks.

**-M, --mail**

:   File(s) contains a list of _mails_ separated by line breaks.

**-U, --url**

:   File(s) contains a list of _urls_ separated by line breaks.

**-I, --ip**

:   File(s) contains a list of _ips_ separated by line breaks.


Required
--------

**--hp-key**=_apikey_

:   _API key_ for **Have I Been Pwned** (hexadecimal format).

**--vt-key**=_apikey_

:   _API key_ for **VirusTotal** (hexadecimal format).

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

ENVIRONMENT
===========

**FOX_HP_KEY**

:   The **Have I Been Pwned** API key can also be set through this environment variable.

**FOX_VT_KEY**

:   The **VirusTotal** API key can also be set through this environment variable.

EXAMPLES
========

fox check sample.exe

:   Check a suspicious file by hash.

fox check -M users.txt

:   Check a list of email addresses.

echo 8.8.8.8 | fox check -I -

:   Check an IP address directly from input.

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
