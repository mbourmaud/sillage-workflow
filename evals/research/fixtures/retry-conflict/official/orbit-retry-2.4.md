# Orbit Retry 2.4 cancellation

`cancelPending()` prevents scheduled attempts from starting. An active attempt
continues until it settles. Applications that require active transport
cancellation must supply an `AbortHandle` and call `abortActive(handle)`.

This page applies to the 2.4 release line.
