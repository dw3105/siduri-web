# W2 · A5 contrast

base: 742e376
lane: A5

## A5 — “the refusal line meets AA contrast”

did: Changed only `.comment-refused`'s text colour from `var(--muted)` to `var(--ink)`, preserving its mixed background, dashed border, monospace face, size, line height, margin and padding, static/comments_a7.css:60.
ran: `node -v && make build && make a8-accessibility`
saw: Before the CSS change, Node was `v20.20.2`; axe-core 4.13.0 reported `1 Accessibility issues detected.` for `.comment-refused` on `/journal/hello-siduri/index.html` and exited 1.
red proof: The original `color: var(--muted)` was the red input. Axe detail output printed light `fgColor: #6e6a60`, `bgColor: #ebe8e1`, `contrastRatio: 4.4`, `expectedContrastRatio: 4.5:1`; forced-dark axe printed `fgColor: #b1aa9e`, `bgColor: #292825`, `contrastRatio: 6.39`, `expectedContrastRatio: 4.5:1`. The light result is below the AA threshold and produced the one violation.
ran: `node -p '"node " + process.versions.node + " (major " + process.versions.node.split(".")[0] + ")"' && make a8-accessibility`
saw: `node 20.20.2 (major 20)`; axe's own output reported `0 violations found!` for every page, including `/journal/hello-siduri/index.html`, followed by `Testing complete of 31 pages`; exit 0.
ran: `set -eu; port=$(python3 tools/free_port.py); (cd dist && python3 -m http.server "$port" --bind 127.0.0.1 >/tmp/siduri-a5-detail-http.log 2>&1) & server=$!; trap 'kill "$server" 2>/dev/null || true' EXIT HUP INT TERM; python3 tools/wait_serving.py "$port"; npx --yes @axe-core/cli@4 --include=.comment-refused --rules=color-contrast --stdout "http://127.0.0.1:$port/journal/hello-siduri/index.html" >/tmp/a5-after-light.json; npx --yes @axe-core/cli@4 --include=.comment-refused --rules=color-contrast --stdout --chrome-options=--force-dark-mode "http://127.0.0.1:$port/journal/hello-siduri/index.html" >/tmp/a5-after-dark.json; node -e 'for (const [name, file] of [["light", "/tmp/a5-after-light.json"], ["dark", "/tmp/a5-after-dark.json"]]) { const result = require(file)[0]; const rule = [...result.violations, ...result.passes].find(item => item.id === "color-contrast"); const data = rule.nodes[0].any[0].data; console.log(name + ": fg=" + data.fgColor + " bg=" + data.bgColor + " ratio=" + data.contrastRatio + " expected=" + data.expectedContrastRatio); }'`
saw: After the change, axe detail output printed light `fgColor: #1d1d1b`, `bgColor: #ebe8e1`, `contrastRatio: 13.79`, `expectedContrastRatio: 4.5:1`; forced-dark axe printed `fgColor: #ece8df`, `bgColor: #292825`, `contrastRatio: 12.05`, `expectedContrastRatio: 4.5:1`. Both dark and light results were measured, not computed.
notes: `go build ./...`, `go test ./internal/site/`, and the prescribed two-build checksum comparison passed. The full `make check` remains for the integrator, as COMMON.md directs lanes not to run the full gate; this lane's `make a8-accessibility` target is green on Node major 20.
