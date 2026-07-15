% FOX HUNT(1) Version 5 | Fox Documentation

NAME
====

**fox** — The Forensic Examiner's Swiss Army Knife

SYNOPSIS
========

| **fox** **hunt** \[_flags_ ...] \[**local** | _paths_ ...]

DESCRIPTION
===========

Hunt suspicious activities by carving events from file(s). Please be aware that, using the **--sort** flag will buffer all found events in memory. For large sets of data this could be very slow and take a serious amount of memory. All timestamps will be normalized to **UTC**.

FLAGS
=====

**-a, --all**

:   Show logs with all severities.

**-s, --sort**

:   Show logs sorted by timestamp.

**-u, --uniq**

:   Show logs that are unique by **XXH3** hash. The calculated hash has 64-bits and is highly unlikely, but still possible, to collide with the another key.

**-j, --json**

:   Show logs as JSON objects.

**-l, --jsonl**

:   Show logs as JSON lines.

**-t, --triage**

:   Show logs in Triage format. Implies **--sort** and **--uniq** flags.

**-p, --parquet**

:   Save logs as Parquet file.

Filter Flags
------------

**-N, --min**=_time_

:   Minimum event _time_ in **RFC3339** format. Example: _2026-12-31T12:00:00.0Z_.

**-X, --max**=_time_

:   Maximum event _time_ in **RFC3339** format. Example: _2026-12-31T12:00:00.0Z_.

**-R, --rule**=_file_

:   Filter using Sigma rules _file_.

Stream Flags
------------

**-U, --url**=_url_

:   Stream events using **CEF** schema to _url_.

**-E, --ecs**=_url_

:   Stream events using **ECS** schema to _url_.

**-H, --hec**=_url_

:   Stream events using **HEC** schema to _url_.

**-A, --auth**=_token_

:   Use auth _token_ with **HEC** streaming. Must be specified without the 'Splunk' prefix.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. If **local** is specified as _path_, a built-in list of known locations will be processed.

EXAMPLES
========

$ fox hunt -t *.dd

:   Hunt down critical events.

$ fox hunt -ap local

:   Save local events as Parquet.

$ fox hunt -E http://127.0.0.1:8080 *.evtx

:   Send events to an Elastic Stack.

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

**fox(1)**, **sort(1)**, **uniq(1)**
