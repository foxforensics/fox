% FOX HUNT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hunt** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Hunt suspicious activities by carving events from file(s). If no path is specified, a built-in list of known locations will be processed.

FLAGS
=====

**-a, --all**

:   Show logs with all severities.

**-s, --sort**

:   Show logs sorted by timestamp (slow).

**-u, --uniq**

:   Show logs that are unique.

**-j, --json**

:   Show logs as JSON objects.

**-J, --jsonl**

:   Show logs as JSON lines.

**-P, --parquet**

:   Save logs as Parquet (very fast).

**-S, --sqlite**

:   Save logs as SQLite3 (very slow).

Hunter Flags
------------

**-b, --block**=_size_

:   Block _size_ for event carving.

Filter Flags
------------

**-R, --rule**=_file_

:   Filter using Sigma Rules _file_ (slow).

**-D, --dist**=_length_

:   Filter using Levenshtein distance (slow).

Stream Flags
------------

**-U, --url**=_server_

:   Stream events to _server_ address.

**-A, --auth**=_token_

:   Stream events using auth _token_.

**-E, --ecs**

:   Use **ECS** schema for streaming.

**-H, --hec**

:   Use **HEC** schema for streaming.

ALIASES
=======

**--logstash**

:   Alias for **-EUhttp://localhost:8080**.

**--splunk**

:   Alias for **-HUhttp://localhost:8088/...**.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**.

EXAMPLES
========

$ fox hunt -u *.dd

:   Hunt down critical events.

$ fox hunt -aP *.evtx

:   Save all events as Parquet.

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

**fox(1)**, **uniq(1)**
