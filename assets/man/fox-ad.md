% FOX AD(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **ad** \[_flags_ ...] **NTDS** **SYSTEM**

DESCRIPTION
===========

Extract **NTLM** password hashes and account infos from Active Directory offline databases.

FLAGS
=====

Account Flags
-------------

**-j, --json**

:   Show accounts as JSON objects.

**-J, --jsonl**

:   Show accounts as JSON lines.

Secrets Flags
-------------

**-H, --history**

:   Extract also the **LM** and **NT** hash history.

**--lm**

:   Extract just the **LM** hashes (hashcat: _3000_).

**--nt**

:   Extract just the **NT** hashes (hashcat: _1000_).

POSITIONAL ARGUMENTS
====================

The Active Directory offline database file followed by the Windows System registry hive.

EXAMPLES
========

$ fox ad -H NTDS.dit SYSTEM

:   Show NTLM secrets.

$ fox ad -j NTDS.dit SYSTEM

:   Show account infos.

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