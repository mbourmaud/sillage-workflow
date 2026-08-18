---
name: frontend-architecture
description: Use automatically for frontend, UI, route, form, state, accessibility, responsive, or interaction work; shape a small vertical with explicit loading, empty, error, success, and real-browser evidence.
metadata:
  namespace: sillage
  qualified-name: "sillage:frontend-architecture"
---

# Build a coherent UI vertical

Start from the user journey and observable states, not from a component list.
Keep framework choices local to the project; the architectural contract is
portable.

## Boundary checklist

- Route/page owns composition and URL semantics; feature hooks own async,
  caching, abort, retry, and mutation policy.
- Domain/application contracts own business meaning; views do not import
  server implementations or duplicate validation rules.
- Components have one visual or interaction responsibility and stay small;
  extract repeated orchestration before a file becomes a mini-application.
- Every meaningful request has loading, empty, success, recoverable error, and
  retry/disabled states. Preserve unsaved input and cancellation behavior.
- Keyboard navigation, labels, focus, contrast, responsive layout, reduced
  motion, and screen-reader semantics are acceptance criteria when relevant.
- Use the project's existing design-system primitives before inventing a
  widget. Test the state contract at the narrowest layer and exercise the real
  browser boundary for layout or interaction risk.

Return:

```text
Journey: <user action → visible result>
States: <loading / empty / success / error / disabled>
Ownership: <route, feature policy, component, contract>
Accessibility: <observable checks>
Browser proof: <route, viewport, screenshot or interaction>
```

If a visual request hides a product decision, route to `sillage:shape`. If a
data contract or domain invariant is wrong, use `sillage:ddd` or a data lens;
do not fix it with view conditionals.
