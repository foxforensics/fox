% FOX DUMP(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **dump** ntds.dit system

DESCRIPTION
===========

Dumps sensitive data from files. This mode enforces the **--no-convert** flag.

FLAGS
=====

**-V, --vss**

:   Dumps data using a Volume Shadow Copy (VSS). Using this flag **WILL ALTER THE FILESYSTEM** and requires a manual confirmation. You have been warned.

**-j, --json**

:   Shows data as JSON objects.

**-J, --jsonl**

:   Shows data as JSON lines.

**--nt**

:   Shows only the NT hashes.

**--lm**

:   Shows only the LM hashes.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox dump ntds.dit system

:   Dump NTLM hashes

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
