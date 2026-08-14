# Security policy

Please report vulnerabilities privately through GitHub's security advisory
feature rather than a public issue.

Sillage's core validator is intentionally read-only by default. Changes that
introduce command execution, credential access, network writes, destructive
filesystem behavior, or automatic publication require explicit threat-model
review and human authorization boundaries.
