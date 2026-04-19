% FOX DUMP(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **dump** **system** \[**ntds**]

DESCRIPTION
===========

Dump critical system secrets from the Windows Registry and Active Directory databases. This command enforces the **--no-convert** flag.

FLAGS
=====

**-j, --json**

:   Dump data as JSON objects.

**-J, --jsonl**

:   Dump data as JSON lines.

Registry Flags
--------------

**-K, --bootkey**

:   Dump the host bootkey.

Active Directory Flags
----------------------

**--only-lm**

:   Extract only the **LM** hashes (hashcat: _3000_).

**--only-nt**

:   Extract only the **NT** hashes (hashcat: _1000_).

POSITIONAL ARGUMENTS
====================

The Windows System registry hive followed by the Active Directory database (optional).

EXAMPLES
========

$ fox dump system -K

:   Dump the BootKey from registry.

$ fox dump system ntds.dit

:   Dump NTLM password hashes.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://foxforensics.dev/fox/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxforensics.dev/fox_>

SEE ALSO
========

**fox(1)**
