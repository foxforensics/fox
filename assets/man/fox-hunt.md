% FOX HUNT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hunt** \[_flags_ ...] \[**local** | _paths_ ...]

DESCRIPTION
===========

Hunt suspicious activities by carving events from file(s). Please be aware that, using the **--sort** flag will buffer all found events in memory. For large sets of data this could be very slow and take a serious amount of memory.

FLAGS
=====

**-a, --all**

:   Show logs with all severities.

**-s, --sort**

:   Show logs sorted by timestamp.

**-u, --uniq**

:   Show logs that are unique.

**-j, --json**

:   Show logs as JSON objects.

**-J, --jsonl**

:   Show logs as JSON lines.

**-p, --parquet**

:   Save logs as Parquet file.

Sigma Flags
-----------

**-R, --rule**=_file_

:   Filter using Sigma rules _file_.

Stream Flags
------------

**-U, --url**=_url_

:   Stream events to a server or broker.

**-E, --ecs**=_url_

:   Stream events with **ECS** schema to _URL_.

**-H, --hec**=_url_

:   Stream events with **HEC** schema to _URL_.

**-A, --auth**=_token_

:   Use auth _token_ for **HEC** streaming. Must be specified without the 'Splunk' prefix.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. If **local** is specified as _path_, a built-in list of known locations will be processed.

EXAMPLES
========

$ fox hunt -u *.dd

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

**fox(1)**, **uniq(1)**
