#!/usr/bin/env python3
"""Print a free loopback port.

The accessibility gate used to hardcode 8765, which is taken on the build host
by an unrelated service. `http.server` failed to bind, the gate did not notice,
and axe scanned **that** application instead -- reporting 122 colour-contrast
violations against a page this repository does not contain. A verdict from the
wrong target is worse than no verdict.
"""
import socket

with socket.socket() as s:
    s.bind(("127.0.0.1", 0))
    print(s.getsockname()[1])
