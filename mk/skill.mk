.PHONY: skill-check skill-install skill-selftest

SKILL_SRC ?= skills/siduri-code
SKILL_LIVE ?= $(HOME)/.claude/skills/siduri-code
SKILL_MANIFEST ?= skills/MANIFEST.tsv
SKILLS_ROOT ?= $(HOME)/.claude/skills

# This is a developer-host check. CI does not have ~/.claude/skills, so it is
# intentionally not a prerequisite of the shared `check` target.
skill-check: skill-selftest
	@python3 -u tools/skill_install.py --src "$(SKILL_SRC)" \
		--live "$(SKILL_LIVE)" --manifest "$(SKILL_MANIFEST)" \
		--skills-root "$(SKILLS_ROOT)" --check

skill-selftest:
	@python3 -u tools/skill_install.py --selftest

skill-install: skill-selftest
	@python3 -u tools/skill_install.py --src "$(SKILL_SRC)" --live "$(SKILL_LIVE)"
