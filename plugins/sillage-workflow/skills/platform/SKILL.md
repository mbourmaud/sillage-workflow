---
name: platform
description: Use automatically when a change depends on desktop, mobile, CLI, service, embedded, offline, process, filesystem, packaging, or operating-system behavior; keep platform policy at an explicit adapter boundary.
metadata:
  namespace: sillage
  qualified-name: "sillage:platform"
---

# Respect the execution environment

Platform architecture is the contract between software and its runtime:
process lifecycle, resources, permissions, storage, input, distribution, and
recovery. Keep product policy portable and isolate OS behavior behind explicit
boundaries.

## Platform pass

1. Identify lifecycle, availability, resource, permission, and distribution
   constraints for the target environment.
2. Model startup, shutdown, interruption, crash recovery, updates, rollback,
   offline behavior, and data migration.
3. For desktop or mobile, consider windows/scenes, focus, keyboard/touch,
   accessibility, deep links, single-instance behavior, local storage,
   sandboxing, notifications, battery, and network transitions when relevant.
4. For CLI or services, consider signals, exit codes, stdin/stdout/stderr,
   cancellation, supervision, logs, configuration, and graceful termination.
5. For embedded or resource-bound targets, make memory, CPU, storage, timing,
   power, and hardware failure budgets explicit.
6. Test the real platform boundary in a representative environment; a mocked
   filesystem or process is not proof of packaging or lifecycle behavior.

Return:

```text
Environment: <platform and lifecycle>
Constraint: <resource, permission, OS, or distribution rule>
Boundary: <portable policy and platform adapter>
Failure/recovery: <interruption, crash, offline, update, or rollback>
Proof: <native, integration, device, browser, process, or packaging check>
```

Do not choose Electron, Swift, .NET, Qt, a mobile toolkit, or a packaging
system from this lens. The project profile and human decision own that choice.
