% FOX HUNT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hunt** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Hunts suspicious activities.

FLAGS
=====

**-a, --all**

:   Shows logs with all severities.

**-s, --sort**

:   Shows logs sorted by timestamp (slow).

**-u, --uniq**

:   Shows logs that are unique.

**-j, --json**

:   Shows logs as JSON objects.

**-J, --jsonl**

:   Shows logs as JSON lines.

**-Q, --sqlite**

:   Saves logs to SQLite3 DB (very slow).

Filter Flags
------------

**-R, --rule**=_file_

:   Filters using Sigma Rules _file_ (slow).

**-D, --dist**=_length_

:   Filters using Levenshtein distance (slow).

Stream Flags
------------

**-U, --url**=_server_

:   Streams events to _server_ address.

**-A, --auth**=_token_

:   Streams events using auth _token_.

**-E, --ecs**

:   Uses **ECS** schema for streaming.

**-H, --hec**

:   Uses **HEC** schema for streaming.

ALIASES
=======

**-L, --logstash**

:   Alias for **-E -Uhttp://localhost:8080**.

**-S, --splunk**

:   Alias for **-H -Uhttp://localhost:8088/...**.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

fox hunt -sv ./**/*.E01

:   Hunt down suspicious events.

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

**fox(1)**, **uniq(1)**
