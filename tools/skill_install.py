#!/usr/bin/env python3
"""Install and check the repository-owned Claude Code skills.

    python3 tools/skill_install.py --src DIR --live DIR
    python3 tools/skill_install.py --src DIR --live DIR --src DIR --live DIR
    python3 tools/skill_install.py --src DIR --live DIR --check
    python3 tools/skill_install.py --selftest

The frontmatter ``name:`` is the skill's identity.  A sibling directory with
that name under a different directory name is a strand: Claude Code can load
it, while a path-only check no longer sees it.  ``--check`` never changes a
directory.  The installer copies the source into an exact, clean live copy and
removes strands for that skill.

The optional manifest check is deliberately separate from content comparison.
It records provenance for skills this repository does not own, but only the
repository-owned source is compared here.

Stdlib only.
"""

from __future__ import annotations

import argparse
import hashlib
import io
import pathlib
import shutil
import sys
import tempfile


FENCE = "---"
MANIFEST_HEADER = ("skill", "canonical", "checked_here")
CHECKED = "yes"
NOT_CHECKED = "not checked here"
IGNORED_DIRS = frozenset({"__pycache__"})


def skill_name(directory: pathlib.Path) -> str | None:
    """Return the frontmatter name in *directory*, or None if it is not one."""
    path = directory / "SKILL.md"
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except (OSError, UnicodeDecodeError):
        return None
    if not lines or lines[0].strip() != FENCE:
        return None
    for line in lines[1:]:
        if line.strip() == FENCE:
            return None
        key, separator, value = line.partition(":")
        if separator and key.strip() == "name":
            return value.strip() or None
    return None


def files(root: pathlib.Path) -> dict[str, str]:
    """Return relative file names and byte digests below *root*."""
    result: dict[str, str] = {}
    for path in sorted(root.rglob("*")):
        relative = path.relative_to(root)
        if IGNORED_DIRS & set(relative.parts):
            continue
        if path.is_file():
            result[relative.as_posix()] = hashlib.sha256(path.read_bytes()).hexdigest()
    return result


def differences(live: pathlib.Path, source: pathlib.Path) -> list[str]:
    """Return relative paths whose membership or bytes differ."""
    live_files, source_files = files(live), files(source)
    return sorted(
        set(live_files) ^ set(source_files)
        | {path for path in set(live_files) & set(source_files) if live_files[path] != source_files[path]}
    )


def strands(live: pathlib.Path, name: str) -> list[pathlib.Path]:
    """Find non-symlink siblings carrying exactly the owned skill *name*."""
    parent = live.parent
    if not parent.is_dir():
        return []
    result = []
    for entry in sorted(parent.iterdir()):
        if entry.name == live.name or entry.is_symlink() or not entry.is_dir():
            continue
        if skill_name(entry) == name:
            result.append(entry)
    return result


def check(source: pathlib.Path, live: pathlib.Path, name: str, out) -> int:
    """Check one source/live pair without changing either directory."""
    stranded = strands(live, name)
    for path in stranded:
        print(
            f"skill-check: {path} carries name '{name}' under another directory "
            "name - Claude Code loads it and this check no longer compares it.",
            file=out,
        )

    failed = bool(stranded)
    if not live.is_dir() or live.is_symlink():
        print(
            f"skill-check: no live copy at {live} - nothing to compare on this machine.",
            file=out,
        )
        failed = True
    else:
        differing = differences(live, source)
        if differing:
            print("skill-check: the live skill has drifted from this repo:", file=out)
            for relative in differing:
                print(f"  differs: {relative}", file=out)
            print(f"  repo ahead of this host:  make skill-install", file=out)
            failed = True
        else:
            print(f"skill-check: live copy matches {source}.", file=out)

    # The warning belongs to every state `skill-install` would act on, not only
    # drift. It first sat inside the drift branch, while the deletion it warns
    # about is in the strand path - a warning placed where it cannot fire in the
    # case it exists for. Found by `siduri-reviewer`, docs/FINDINGS.md row 0093.
    if failed:
        print(
            "  before running `make skill-install`: docs/FINDINGS.md row 0093.\n"
            "  install() removes a stranded directory under the live skills path,\n"
            "  which holds skills this repo does not own, and has never run\n"
            "  against it.",
            file=out,
        )

    if stranded:
        print("  remove the strand:        make skill-install", file=out)
    return 1 if failed else 0


def install(source: pathlib.Path, live: pathlib.Path, name: str, out) -> int:
    """Replace the live copy with an exact source copy and prune its strands."""
    if live.is_symlink() or (live.exists() and not live.is_dir()):
        print(f"skill-install: live path is not a directory: {live}", file=sys.stderr)
        return 1

    live.parent.mkdir(parents=True, exist_ok=True)
    # Stage before replacing the old directory.  A failed copy leaves the old
    # live copy available, and an interrupted strand prune can only leave the
    # new live copy plus a strand, which the next check reports.
    with tempfile.TemporaryDirectory(prefix=f".{live.name}-", dir=live.parent) as directory:
        staged = pathlib.Path(directory) / live.name
        shutil.copytree(source, staged)
        if live.exists():
            shutil.rmtree(live)
        staged.rename(live)

    for path in strands(live, name):
        shutil.rmtree(path)
        print(f"skill-install: removed {path}, which carried name '{name}'.", file=out)
    print(f"skill-install: {live} now matches {source}.", file=out)
    return 0


def install_many(
    pairs: list[tuple[pathlib.Path, pathlib.Path, str]], out
) -> int:
    """Install every source/live pair, retaining exact-name strand scoping."""
    status = 0
    for source, live, name in pairs:
        status = max(status, install(source, live, name, out))
    return status


def manifest_rows(path: pathlib.Path) -> dict[str, tuple[str, str]]:
    """Read the three-column provenance manifest, rejecting malformed rows."""
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except OSError as error:
        raise ValueError(f"cannot read {path}: {error}") from error

    rows: dict[str, tuple[str, str]] = {}
    header_seen = False
    for line_number, line in enumerate(lines, 1):
        if not line.strip() or line.lstrip().startswith("#"):
            continue
        fields = tuple(field.strip() for field in line.split("\t"))
        if not header_seen:
            if fields != MANIFEST_HEADER:
                raise ValueError(
                    f"{path}:{line_number}: expected tab-separated header "
                    "skill<TAB>canonical<TAB>checked_here"
                )
            header_seen = True
            continue
        if len(fields) != 3 or not all(fields):
            raise ValueError(f"{path}:{line_number}: expected three non-empty fields")
        name, canonical, checked = fields
        if checked not in {CHECKED, NOT_CHECKED}:
            raise ValueError(
                f"{path}:{line_number}: checked_here must be '{CHECKED}' or '{NOT_CHECKED}'"
            )
        if name in rows:
            raise ValueError(f"{path}:{line_number}: duplicate skill row {name!r}")
        rows[name] = (canonical, checked)
    if not header_seen:
        raise ValueError(f"{path}: manifest has no header")
    return rows


def deployed_directories(root: pathlib.Path) -> list[pathlib.Path]:
    """Return every direct child directory Claude Code could load."""
    if not root.is_dir() or root.is_symlink():
        return []
    return sorted(
        entry
        for entry in root.iterdir()
        if entry.is_dir() or entry.is_symlink()
    )


def check_manifest(
    manifest: pathlib.Path,
    skills_root: pathlib.Path,
    source_names: list[str],
    out,
) -> int:
    """Check manifest coverage and the in-repository checked rows."""
    try:
        rows = manifest_rows(manifest)
    except ValueError as error:
        print(f"skill-manifest: {error}", file=out)
        return 1

    deployed = deployed_directories(skills_root)
    if not skills_root.is_dir() or skills_root.is_symlink():
        print(f"skill-manifest: no deployed skills directory at {skills_root}.", file=out)
        return 1

    deployed_names = {path.name for path in deployed}
    manifest_names = set(rows)
    missing = sorted(deployed_names - manifest_names)
    # A checked-here source is recorded before the integrator deploys it to
    # this host, so its repository directory is also an expected location.
    stale = sorted(manifest_names - deployed_names - set(source_names))
    checked = sorted(name for name, (_, value) in rows.items() if value == CHECKED)
    failed = False
    if missing:
        failed = True
        for name in missing:
            print(f"skill-manifest: deployed directory has no row: {name}", file=out)
    if stale:
        failed = True
        for name in stale:
            print(f"skill-manifest: row has no deployed directory: {name}", file=out)
    expected = sorted(source_names)
    if checked != expected:
        failed = True
        print(
            f"skill-manifest: checked_here=yes rows are {checked!r}; expected {expected!r}",
            file=out,
        )
    if failed:
        return 1
    return 0


def _remove(path: pathlib.Path) -> None:
    if path.is_symlink() or path.is_file():
        path.unlink()
    elif path.exists():
        shutil.rmtree(path)


def _fixture(source: pathlib.Path, name: str = "siduri-code") -> None:
    source.mkdir(parents=True)
    (source / "SKILL.md").write_text(
        f"---\nname: {name}\ndescription: selftest\n---\n", encoding="utf-8"
    )
    (source / "payload.txt").write_text("canonical\n", encoding="utf-8")


def _captured(function, *args) -> tuple[int, str]:
    output = io.StringIO()
    return function(*args, output), output.getvalue()


def selftest() -> int:
    """Exercise every red condition and its adjacent passing fixture."""
    with tempfile.TemporaryDirectory(prefix="siduri-skill-source-") as source_dir:
        with tempfile.TemporaryDirectory(prefix="siduri-skill-live-") as live_dir:
            source = pathlib.Path(source_dir) / "siduri-code"
            live = pathlib.Path(live_dir) / "siduri-code"
            _fixture(source)

            shutil.copytree(source, live)
            (live / "payload.txt").write_text("live edit\n", encoding="utf-8")
            status, output = _captured(check, source, live, "siduri-code")
            assert status == 1 and "differs: payload.txt" in output
            print("skill-install selftest: differs red — two temporary copies, live payload.txt changed")
            assert install(source, live, "siduri-code", io.StringIO()) == 0
            status, output = _captured(check, source, live, "siduri-code")
            assert status == 0 and "matches" in output
            print("skill-install selftest: differs green — installer restored exact canonical bytes")

            _remove(live)
            status, output = _captured(check, source, live, "siduri-code")
            assert status == 1 and "no live copy" in output
            print("skill-install selftest: absent red — --live points at missing temporary directory")
            assert install(source, live, "siduri-code", io.StringIO()) == 0
            status, output = _captured(check, source, live, "siduri-code")
            assert status == 0
            print("skill-install selftest: absent green — installer created the live directory")

            strand = live.parent / "renamed-siduri"
            shutil.copytree(source, strand)
            status, output = _captured(check, source, live, "siduri-code")
            assert status == 1 and "carries name 'siduri-code'" in output
            print("skill-install selftest: strand red — sibling renamed-siduri carries frontmatter name")
            assert install(source, live, "siduri-code", io.StringIO()) == 0
            status, output = _captured(check, source, live, "siduri-code")
            assert status == 0 and not strand.exists()
            print("skill-install selftest: strand green — installer removed the sibling strand")

            skills_root = pathlib.Path(live_dir) / "deployed"
            skills_root.mkdir()
            (skills_root / "siduri-code").mkdir()
            (skills_root / "foreign-skill").mkdir()
            manifest = pathlib.Path(live_dir) / "MANIFEST.tsv"
            manifest.write_text(
                "skill\tcanonical\tchecked_here\n"
                "siduri-code\tskills/siduri-code\tyes\n",
                encoding="utf-8",
            )
            status, output = _captured(
                check_manifest, manifest, skills_root, ["siduri-code"]
            )
            assert status == 1 and "foreign-skill" in output
            print("skill-manifest selftest: completeness red — deployed foreign-skill has no row")
            with manifest.open("a", encoding="utf-8") as stream:
                stream.write("foreign-skill\texternal:fixture\tnot checked here\n")
            status, output = _captured(
                check_manifest, manifest, skills_root, ["siduri-code"]
            )
            assert status == 0 and output == ""
            print("skill-manifest selftest: completeness green — every deployed directory is listed")

            with manifest.open("a", encoding="utf-8") as stream:
                stream.write("siduri-contract\tsiduri-web:skills/siduri-contract\tyes\n")
            status, output = _captured(
                check_manifest,
                manifest,
                skills_root,
                ["siduri-code", "siduri-contract"],
            )
            assert status == 0 and output == ""
            print(
                "skill-manifest selftest: repo-owned green — an undeployed "
                "checked-here source is expected before integration"
            )

            with manifest.open("a", encoding="utf-8") as stream:
                stream.write("stale-skill\texternal:fixture\tnot checked here\n")
            status, output = _captured(
                check_manifest,
                manifest,
                skills_root,
                ["siduri-code", "siduri-contract"],
            )
            assert status == 1 and "stale-skill" in output
            print("skill-manifest selftest: stale red — unowned row has no deployed directory")

            multi_root = pathlib.Path(live_dir) / "multi"
            multi_root.mkdir()
            foreign = multi_root / "foreign-hitl-code"
            foreign.mkdir()
            (foreign / "SKILL.md").write_text(
                "---\nname: hitl-code\ndescription: foreign\n---\n",
                encoding="utf-8",
            )
            (foreign / "payload.txt").write_text("foreign\n", encoding="utf-8")
            foreign_before = files(foreign)

            names = (
                "siduri-code",
                "siduri-contract",
                "siduri-pending-tasks",
                "siduri-pending-proving",
            )
            sources = [source]
            for name in names[1:]:
                fixture = pathlib.Path(source_dir) / name
                _fixture(fixture, name)
                sources.append(fixture)
            pairs = [
                (skill_source, multi_root / name, name)
                for skill_source, name in zip(sources, names)
            ]
            assert install_many(pairs, io.StringIO()) == 0
            assert foreign.is_dir() and files(foreign) == foreign_before
            print(
                "skill-install selftest: foreign green — four installs left "
                "foreign hitl-code byte-identical"
            )
    print("skill-install selftest: 11 named red/green cases pass")
    return 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    parser.add_argument("--src", type=pathlib.Path, action="append")
    parser.add_argument("--live", type=pathlib.Path, action="append")
    parser.add_argument("--manifest", type=pathlib.Path)
    parser.add_argument("--skills-root", type=pathlib.Path)
    parser.add_argument("--check", action="store_true", help="report only; change nothing")
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args(argv)
    if args.selftest:
        return selftest()
    if not args.src or not args.live:
        parser.error("--src and --live are required unless --selftest is used")
    if len(args.src) != len(args.live):
        parser.error("each --src must have a matching --live")
    if args.manifest is not None and args.skills_root is None:
        parser.error("--manifest requires --skills-root")
    if args.skills_root is not None and args.manifest is None:
        parser.error("--skills-root requires --manifest")

    pairs: list[tuple[pathlib.Path, pathlib.Path, str]] = []
    for source, live in zip(args.src, args.live):
        if not source.is_dir():
            print(f"skill: no source directory at {source}", file=sys.stderr)
            return 1
        name = skill_name(source)
        if not name:
            print(
                f"skill: {source / 'SKILL.md'} declares no frontmatter name",
                file=sys.stderr,
            )
            return 1
        pairs.append((source, live, name))

    if args.check:
        status = 0
        if args.manifest is not None and args.skills_root is not None:
            status = check_manifest(
                args.manifest,
                args.skills_root,
                [name for _, _, name in pairs],
                sys.stdout,
            )
        for source, live, name in pairs:
            status = max(status, check(source, live, name, sys.stdout))
        return status
    return install_many(pairs, sys.stdout)


if __name__ == "__main__":
    sys.exit(main())
