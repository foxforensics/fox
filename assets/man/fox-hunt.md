% FOX HUNT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hunt** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Hunts suspicious activities by carving events from file(s). This mode enforces the **--no-extract**, **--no-deflate** and **--no-convert** flags unless **--no-strict** is specified.

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

**-P, --parquet**

:   Saves logs as Parquet (very fast).

**-Q, --sqlite**

:   Saves logs as SQLite3 (very slow).

Hunter Flags
------------

**-b, --block**=_size_

:   Block _size_ for event carving.

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

Agent Flags
-----------
**-M, --mcp**

:   Starts as **MCP** server (blocking).

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

fox hunt -u *.dd

:   Hunt down suspicious events.

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

**fox(1)**, **uniq(1)**
