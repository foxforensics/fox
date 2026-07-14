% FOX AD(1) Version 5 | Fox Documentation

NAME
====

**fox** — The Forensic Examiner's Swiss Army Knife

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

:   Lookup hashes using the built-in wordlist. **NT** and **LM** hashes will be replaced in place.

**-h, --history**

:   Extract also the users hash history. Lookup of these hashes is also possible.

**--lm-only**

:   Extract only the **LM** hashes (Hashcat mode _3000_). Output will always be plaintext.

**--nt-only**

:   Extract only the **NT** hashes (Hashcat mode _1000_). Output will always be plaintext.

POSITIONAL ARGUMENTS
====================

The Active Directory offline database file followed by the Windows system registry hive. To refer to paths inside archives, use the archive!file notation.

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
<_https://foxforensics.eu/issues_>

WWW
===

Please visit the project's homepage at:
<_https://foxforensics.eu_>

SEE ALSO
========

**fox(1)**, **impacket-secretsdump(1)**