package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadContentRejectsMissingPlainSummary(t *testing.T) {
	root := fixtureRoot(t, `---
title: Missing summary
slug: missing-summary
date: 2026-08-25
summary: A summary.
tags:
  - method
draft: false
---

Body.`)
	_, err := LoadContent(root)
	if err == nil || !strings.Contains(err.Error(), `missing required frontmatter field "plain_summary"`) {
		t.Fatalf("expected missing plain_summary error, got %v", err)
	}
}

func TestLoadContentRejectsUnknownTag(t *testing.T) {
	root := fixtureRoot(t, `---
title: Unknown tag
slug: unknown-tag
date: 2026-08-25
summary: A summary.
plain_summary: A plain summary.
tags:
  - nonsense
draft: false
---

Body.`)
	_, err := LoadContent(root)
	if err == nil || !strings.Contains(err.Error(), `unknown tag "nonsense"`) {
		t.Fatalf("expected unknown tag error, got %v", err)
	}
}

func fixtureRoot(t *testing.T, post string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "content", "posts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "content", "tags.yml"), []byte("tags:\n  - build-log\n  - tool-release\n  - dogfooding\n  - method\n  - outcome\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "content", "posts", "fixture.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
