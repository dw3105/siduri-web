# W1 · G1 — build integrity

wave: W1
tree: 7667b8d5c067d047bffa402cdb58fc92550a3ccc

## Criterion 1 — A clean checkout builds to `dist/` with one command, no network, no manual steps.

verdict: failed

ran: `clone_dir=$(mktemp -d /tmp/siduri-w1-g1-clone.XXXXXX); git clone --no-local "$PWD" "$clone_dir"; echo "clone=$clone_dir"; echo "tree=$(git -C "$clone_dir" rev-parse HEAD)"; echo "tracked_templ=$(git -C "$clone_dir" ls-files -- '*templ' '*templ.exe' | wc -l)"; test -d "$clone_dir/vendor" && echo 'vendor=present' || echo 'vendor=absent'; grep -Eq '^[[:space:]]*tool([[:space:]]|\(|$)' "$clone_dir/go.mod" && echo 'go_mod_tool=present' || echo 'go_mod_tool=absent'; date -u +%Y-%m-%dT%H:%M:%SZ; set +e; env GOPATH="$clone_dir/.gopath" unshare -rn make -C "$clone_dir" build; status=$?; set -e; echo "build_exit=$status"; test "$status" -ne 0`
at: 2026-08-27T18:19:11Z
saw: clone `/tmp/siduri-w1-g1-clone.jWZOLI` at tree `7667b8d5c067d047bffa402cdb58fc92550a3ccc`; `vendor=absent`; `go_mod_tool=absent`; `/bin/sh: 1: /tmp/siduri-w1-g1-clone.jWZOLI/.gopath/bin/templ: not found`; `make: *** [Makefile:21: build] Error 127`; `build_exit=2`.

- **ran** — `clone_dir=/tmp/siduri-w1-g1-clone.jWZOLI; date -u +%Y-%m-%dT%H:%M:%SZ; test -x "$clone_dir/.gopath/bin/templ" && echo 'clone_templ=present' || echo 'clone_templ=absent'; git -C "$clone_dir" ls-files | rg '(^|/)templ($|\.)' || true; test -d "$clone_dir/vendor" && echo 'clone_vendor=present' || echo 'clone_vendor=absent'; grep -Eq '^[[:space:]]*tool([[:space:]]|\(|$)' "$clone_dir/go.mod" && echo 'clone_go_mod_tool=present' || echo 'clone_go_mod_tool=absent'`
- **at** — 2026-08-27T18:23:50Z
- **saw** — `clone_templ=absent`; no tracked `templ` executable; `clone_vendor=absent`; `clone_go_mod_tool=absent`.

- **ran** — `clone_dir=/tmp/siduri-w1-g1-clone.jWZOLI; date -u +%Y-%m-%dT%H:%M:%SZ; set +e; env GOPATH="$clone_dir/.gopath" GOMODCACHE="$clone_dir/.gomodcache" unshare -rn go list -m all; status=$?; set -e; echo "module_list_exit=$status"; exit 0`
- **at** — 2026-08-27T18:24:13Z
- **saw** — the unnetworked fresh module cache tried `https://proxy.golang.org/github.com/a-h/templ/@v/v0.3.1020.info` and `.mod`, then reported `network is unreachable`; `module_list_exit=1`.

red proof: `unshare -rn` was usable. The clean clone failed before Go compilation because the build recipe directly invokes the absent `$GOPATH/bin/templ`. The missing executable is the immediate blocker; CI supplies it with a separate manual `go install` step. Satisfying the stated criterion is therefore not just a rendered-output change or a one-line `go.mod` annotation: the tool and the Go module source must be available to a fresh checkout without network access or manual provisioning.

- **ran** — `date -u +%Y-%m-%dT%H:%M:%SZ; set +e; timeout --signal=TERM --kill-after=3s 20s unshare -rn env A8_EXTERNAL_LINKS=0 make check; status=$?; set -e; echo "make_check_exit=$status"; exit 0`
- **at** — 2026-08-27T18:21:49Z
- **saw** — build, pre-P3, HTML structure, budgets, headers, and internal-only link checks ran; linkcheck saw `518 internal links, 4 external links, 0 failures, internal only; external checks skipped (use --external)`; the command then remained in the Node gate and timed out with `make_check_exit=124`.

- **ran** — `date -u +%Y-%m-%dT%H:%M:%SZ; make -n check | rg -n 'npx|linkcheck|go install|curl|wget'`
- **at** — 2026-08-27T18:23:50Z
- **saw** — `a8-html` runs `npx --yes html-validate@8`; `a8-accessibility` runs `npx --yes @axe-core/cli@4`; with `A8_EXTERNAL_LINKS=1`, `a8-links` runs `tools/linkcheck.py --external`, whose external branch performs HTTP HEAD/GET requests. The no-network check used `A8_EXTERNAL_LINKS=0`, so its external links were intentionally skipped. The two `npx` package resolutions are the network-sensitive check steps.

notes: The criterion is failed on a clean clone. The temporary clone was made with `git clone --no-local` from this worktree, not copied, so ignored and untracked files did not supply prerequisites. `go.mod` requires `github.com/a-h/templ v0.3.1020`, but the checkout has no vendor directory and no `tool` directive; a fresh module cache consequently also needs network access. `make check` cannot be credited as a no-network command: the bounded run reached the local gates but stalled at the first Node package resolution, and the accessibility gate would resolve its second package on a Node 20 runner. CI additionally opts into external HTTP link checks.

## Criterion 2 — Two consecutive builds from identical input produce byte-identical output.

verdict: held

ran: `first_dir=$(mktemp -d /tmp/siduri-w1-g1-build1.XXXXXX); second_dir=$(mktemp -d /tmp/siduri-w1-g1-build2.XXXXXX); first_manifest="$first_dir.manifest"; second_manifest="$second_dir.manifest"; date -u +%Y-%m-%dT%H:%M:%SZ; PATH="$(go env GOPATH)/bin:$PATH" templ generate; go run ./cmd/siduri build --output "$first_dir"; go run ./cmd/siduri build --output "$second_dir"; (cd "$first_dir" && find . -type f -print0 | sort -z | xargs -0 sha256sum) > "$first_manifest"; (cd "$second_dir" && find . -type f -print0 | sort -z | xargs -0 sha256sum) > "$second_manifest"; echo "first_files=$(wc -l < "$first_manifest") second_files=$(wc -l < "$second_manifest")"; if cmp -s "$first_manifest" "$second_manifest"; then echo 'recursive_sha256=identical'; else echo 'recursive_sha256=different'; exit 1; fi; echo "first_dir=$first_dir"; echo "second_dir=$second_dir"`
at: 2026-08-27T18:22:44Z
saw: `first_files=41 second_files=41`; `recursive_sha256=identical` for `/tmp/siduri-w1-g1-build1.MaqACZ` and `/tmp/siduri-w1-g1-build2.16SzVf`.

- **ran** — `date -u +%Y-%m-%dT%H:%M:%SZ; unshare -rn go test ./internal/site -run '^TestA7RepositoryCommentsRenderAsOneLevelThread$' -count=1`
- **at** — 2026-08-27T18:22:48Z
- **saw** — `ok github.com/dw3105/siduri-web/internal/site 0.065s`.

red proof: A scratch copy was patched only with `time.Now().UTC().Format(time.RFC3339Nano)` appended to its generated stylesheet. The two-build recursive checksum comparison then returned `recursive_sha256_cmp_exit=1`; the diff showed `./site.css` changing from `556b4224395888bbc4005907e93e688de28eaa9fc044928a5f05966f75d12b6e` to `44747385d7e24df41487c9d14957b99cdee997e0c682d70d048a52de64d4e570`. The scratch copy and probe outputs were temporary and moved to the desktop trash after evidence collection.

- **ran** — `scratch_dir=/tmp/siduri-w1-g1-nondeterministic.59Jgqt; first_dir=$(mktemp -d /tmp/siduri-w1-g1-red1.XXXXXX); second_dir=$(mktemp -d /tmp/siduri-w1-g1-red2.XXXXXX); first_manifest="$first_dir.manifest"; second_manifest="$second_dir.manifest"; date -u +%Y-%m-%dT%H:%M:%SZ; (cd "$scratch_dir" && unshare -rn go run ./cmd/siduri build --output "$first_dir"); (cd "$scratch_dir" && unshare -rn go run ./cmd/siduri build --output "$second_dir"); (cd "$first_dir" && find . -type f -print0 | sort -z | xargs -0 sha256sum) > "$first_manifest"; (cd "$second_dir" && find . -type f -print0 | sort -z | xargs -0 sha256sum) > "$second_manifest"; echo "first_files=$(wc -l < "$first_manifest") second_files=$(wc -l < "$second_manifest")"; set +e; cmp -s "$first_manifest" "$second_manifest"; cmp_status=$?; diff -u "$first_manifest" "$second_manifest" | sed -n '1,12p'; set -e; echo "recursive_sha256_cmp_exit=$cmp_status"; test "$cmp_status" -eq 1`
- **at** — 2026-08-27T18:23:17Z
- **saw** — `first_files=41 second_files=41`; recursive checksum comparison failed only on `./site.css`; `recursive_sha256_cmp_exit=1`.

notes: The measured build used `go run ./cmd/siduri build` twice with the repository's published content and compared a sorted SHA-256 record for every output file and relative path. The existing full-tree Go test is `TestA7RepositoryCommentsRenderAsOneLevelThread`; it calls `Build` twice on repository content and compares every file with `a7BuildFiles`, then checks the comment page. This evidence says nothing about different content, and nothing about `.dev-dist` versus `dist`.
