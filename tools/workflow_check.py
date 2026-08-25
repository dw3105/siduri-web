#!/usr/bin/env python3
"""Every workflow file parses as YAML, and there is at least one.

Was `npx -p yaml node -e ...`, which resolved the package on Node 18 and threw
MODULE_NOT_FOUND on Node 20 -- a gate that depended on a Node version nothing
pinned. Python parses YAML here and in every other gate in this file.
"""
import pathlib, sys

try:
    import yaml
except ImportError:
    sys.exit("workflow_check: PyYAML is not installed; the gate cannot run")

d = pathlib.Path(".github/workflows")
files = sorted(f for f in d.glob("*.y*ml"))
if not files:
    sys.exit("workflow_check: no workflow files")
for f in files:
    try:
        yaml.safe_load(f.read_text())
    except yaml.YAMLError as exc:
        sys.exit(f"workflow_check: {f} is not valid YAML: {exc}")
print(f"workflow YAML: valid ({len(files)} file(s))")
