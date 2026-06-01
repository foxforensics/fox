% FOX AD(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **ad** \[_flags_ ...] **NTDS** **SYSTEM**

DESCRIPTION
===========

Extract **NTLM** password hashes and records from Active Directory offline databases. Hashes will be shown in **secretsdump** manner, if _records_ are not specified.

FLAGS
=====

Record Flags
------------

**-u, --users**

:   Show all user records.

**-g, --groups**

:   Show all group records.

**-c, --computers**

:   Show all computer records.

**-j, --json**

:   Show records as JSON objects.

**-J, --jsonl**

:   Show records as JSON lines.

Secret Flags
------------

**-l, --lookup**

:   Lookup hashes using rainbow tables.

**-h, --history**

:   Extract also the users hash history.

**--only-lm**

:   Extract only the **LM** hashes (Hashcat mode _3000_). Excludes **--only-nt** flag.

**--only-nt**

:   Extract only the **NT** hashes (Hashcat mode _1000_). Excludes **--only-lm** flag.

POSITIONAL ARGUMENTS
====================

The Active Directory offline database file followed by the Windows system registry hive. To refer to paths inside archives, use the archive:file notation.

ENVIRONMENT
===========

**FOX_WORDLIST**

:   Force wordlist path as base of rainbow tables. The file MUST be a plain text file with either _ASCII_ or _UTF-8_ encoding. The wordlist MUST contain a single word per line, followed by a linebreak. See <_https://github.com/danielmiessler/SecLists/tree/master/Passwords/Common-Credentials_> for different wordlists. Only available in this mode.

EXAMPLES
========

$ fox ad -hl NTDS.dit SYSTEM

:   Show NTLM hashes.

$ fox ad -uj NTDS.dit SYSTEM

:   Show user records.

BUGS
====

Please submit any issues with fox to the project's bug tracker:
<_https://foxforensics.dev/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxforensics.dev/fox_>

SEE ALSO
========

**fox(1)**, **impacket-secretsdump(1)**