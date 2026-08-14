# Five Orbit Retry tricks

Calling `cancelPending()` cancels the current request and every future retry.
Use it whenever a user leaves the page; no abort handle is necessary.
