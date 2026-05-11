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

**-j, --json**

:   Show accounts as JSON objects.

**-J, --jsonl**

:   Show accounts as JSON lines.

Secrets Flags
-------------

**-L, --lookup**

:   Lookup hashes in the rainbow tables.

**-H, --history**

:   Extract also the users hash history.

**--only-lm**

:   Extract only the **LM** hashes (hashcat type _3000_). Excludes **--only-nt** flag.

**--only-nt**

:   Extract only the **NT** hashes (hashcat type _1000_). Excludes **--only-lm** flag.

POSITIONAL ARGUMENTS
====================

The Active Directory offline database file followed by the Windows system registry hive.

ENVIRONMENT
===========

**FOX_WORDLIST**

:   Force wordlist path as base of rainbow tables. Only available in AD mode.

EXAMPLES
========

$ fox ad -LH NTDS.dit SYSTEM

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

**fox(1)**, **impacket-secretsdump(1)**