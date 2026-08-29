.PHONY: skill-check skill-install skill-selftest

SKILL_NAMES := siduri-code siduri-contract
SKILL_MANIFEST ?= skills/MANIFEST.tsv
SKILLS_ROOT ?= $(HOME)/.claude/skills
SKILL_ARGS = $(foreach skill,$(SKILL_NAMES),--src "skills/$(skill)" --live "$(SKILLS_ROOT)/$(skill)")

# This is a developer-host check. CI does not have ~/.claude/skills, so it is
# intentionally not a prerequisite of the shared `check` target.
skill-check: skill-selftest
	@python3 -u tools/skill_install.py $(SKILL_ARGS) \
		--manifest "$(SKILL_MANIFEST)" --skills-root "$(SKILLS_ROOT)" --check

skill-selftest:
	@python3 -u tools/skill_install.py --selftest

skill-install: skill-selftest
	@python3 -u tools/skill_install.py $(SKILL_ARGS)
