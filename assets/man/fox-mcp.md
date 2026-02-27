% FOX MCP(1) Version 4 | Fox Documentation

NAME
====

**fox** — The Forensic Examiners Swiss Army Knife

SYNOPSIS
========

| **fox** **mcp** \[_flags_ ...] \[_port_]

DESCRIPTION
===========

Init MCP server and block until canceling. Once started, the server will be available under <_http://localhost:3001/mcp_>.

POSITIONAL ARGUMENTS
====================

The port number (default: **3001**).

EXAMPLES
========

fox mcp 8080

:   Init the MCP server

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

**fox(1)**
