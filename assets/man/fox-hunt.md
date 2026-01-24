% FOX(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hunt** \[_flags_ ...] \[_paths_ ...]

DESCRIPTION
===========

Hunt suspicious activities.

FLAGS
=====

**-a, --all**

:   Show logs with all severities.

**-s, --sort**

:   Show logs sorted by timestamp (slow).

**-j, --json**

:   Show logs as JSON objects.

**-J, --jsonl**

:   Show logs as JSON lines.

**-D, --sqlite**

:   Save logs to SQLite3 DB (very slow).

Rule Flags
----------

**-R, --rule**=_file_

:   Filter using a Sigma rule (slow).

Stream Flags
------------

**-U, --url**=_server_

:   Stream events to _server_ address.

**-T, --auth**=_token_

:   Stream events using auth _token_.

**-E, --ecs**

:   Use **ECS** schema for streaming.

**-H, --hec**

:   Use **HEC** schema for streaming.

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

SEE ALSO
========

**fox(1)**
