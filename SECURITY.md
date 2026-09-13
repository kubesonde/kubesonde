# Security Policy

## Reporting a vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, use GitHub's private vulnerability reporting:

1. Go to the repository's **Security** tab.
2. Click **Report a vulnerability**.
3. Fill in the advisory form with as much detail as you can — affected
   version(s), a description of the issue, and steps to reproduce or a
   proof-of-concept if you have one.

This opens a private advisory visible only to you and the maintainers. We will
acknowledge your report, investigate, and keep you updated on the fix and
disclosure timeline.

## Scope

Kubesonde is a **diagnostic and auditing tool**: it observes and reports
connectivity and does not block or modify traffic. Reports we are particularly
interested in include (but are not limited to):

- Privilege escalation or container escape via the injected `debugger` /
  `monitor` ephemeral containers.
- Exposure of sensitive cluster data through the controller HTTP API
  (`:2709`) or the frontend.
- Ability for a probed workload to affect the controller or other namespaces.

## Supported versions

Security fixes are provided for the **latest released version**. Please upgrade
to the latest [release](https://github.com/kubesonde/kubesonde/releases) before
reporting, in case the issue is already fixed.
