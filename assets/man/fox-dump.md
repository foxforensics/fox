% FOX DUMP(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **dump** **system** \[**ntds.dit**]

DESCRIPTION
===========

Dumps sensitive data from files. This mode enforces the **--no-convert** flag.

FLAGS
=====

**-j, --json**

:   Shows data as JSON objects.

**-J, --jsonl**

:   Shows data as JSON lines.

Registry Flags
--------------

**-K, --bootkey**

:   Extracts only the bootkey.

Active Directory Flags
----------------------

**--nt**

:   Extracts only the **NT** hashes.

**--lm**

:   Extracts only the **LM** hashes.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox dump system ntds.dit

:   Dump NTLM hashes

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
