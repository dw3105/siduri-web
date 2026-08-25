#!/usr/bin/env python3
"""Refuse until dist/ is actually being served on the given port.

Binding is not serving, and the previous gate assumed the first implied the
second. It asks for a file only this build produces, so another service holding
the port fails the check rather than passing it.
"""
import sys, time, urllib.request

port = sys.argv[1]
deadline = time.monotonic() + 20
while time.monotonic() < deadline:
    try:
        with urllib.request.urlopen(f"http://127.0.0.1:{port}/index.html", timeout=1) as r:
            if r.status == 200 and b"Siduri" in r.read(4096):
                print(f"wait_serving: dist/ confirmed on {port}")
                sys.exit(0)
    except Exception:
        time.sleep(0.3)
sys.exit(1)
