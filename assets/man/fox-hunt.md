% FOX HUNT(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **hunt** \[_flags_ ...] \[**local** | _paths_ ...]

DESCRIPTION
===========

Hunt suspicious activities by carving events from file(s).

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

**-p, --parquet**

:   Save logs as Parquet file.

Block Flags
-----------

**-B, --block**=_size_

:   Block _size_ for event carving. The _size_ can be either defined as raw bytes or with a size suffix.

Filter Flags
------------

**-R, --rule**=_file_

:   Filter using Sigma Rules _file_.

**-D, --dist**=_length_

:   Filter using Levenshtein distance (slow).

Stream Flags
------------

**-U, --url**=_server_

:   Stream events to a server or broker. Streaming via HTTP(S) is set as the default. To stream via the MQTT V5 protocol, use the **--mqtt** flag and specify a topic. 

**-A, --auth**=_token_

:   Authentication _token_ used for Splunk servers. Must be specified without the 'Splunk' prefix.

**-M, --mqtt**=_topic_

:   Use the MQTT protocol V5 for streaming. Currently only streaming via **TCP** is supported. A _topic_ is required.

Schema Flags
------------

**-e, --ecs**

:   Use ECS schema while streaming.

**-h, --hec**

:   Use HEC schema while streaming.

ALIASES
=======

**--elastic**

:   Alias for **--ecs --url http://localhost:8080**.

**--splunk**

:   Alias for **--hec --url http://localhost:8088/...**.

POSITIONAL ARGUMENTS
====================

Globbing paths to open or '-' to read from **STDIN(4)**. If **local** is specified as _path_, a built-in list of known locations will be processed.

ENVIRONMENT
===========

**FOX_MQTT_QOS**

:   The MQTT protocol Quality of Service level. Range 0 to 2 (default: 1).

**FOX_MQTT_USERNAME**

:   The username used to connect to the given MQTT broker.

**FOX_MQTT_PASSWORD**

:   The password used to connect to the given MQTT broker.

EXAMPLES
========

$ fox hunt -u *.dd

:   Hunt down critical events.

$ fox hunt -ap *.evtx

:   Save all events as Parquet.

$ fox hunt -U http://127.0.0.1:8080 local

:   Send local events to a server.

$ fox hunt -U tcp://127.0.0.1:1883 -M events local

:   Send local events to a broker.

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

**fox(1)**, **uniq(1)**
