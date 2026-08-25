package site

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestC3ToolWithoutRepositoryHasNoRepositoryArtefact(t *testing.T) {
	tool := Tool{
		Title: "Unpublished", Slug: "unpublished", Date: "2026-08-25", Summary: "A tool without a public repository.",
		Language: "Go", Status: "active", Install: "from a local checkout", Example: "unpublished run", Why: "The repository is not public yet.",
	}
	output := renderComponent(t, toolPage(tool, nil, PageData{}, nil))
	for _, forbidden := range []string{"Repository", "None", `<a href="">`, "<a>"} {
		if bytes.Contains(output, []byte(forbidden)) {
			t.Fatalf("tool without a repository contains %q in rendered HTML: %s", forbidden, output)
		}
	}
}

func TestC3InvalidRepositoryNamesToolAndValue(t *testing.T) {
	dir := writeC3Tool(t, `repository: not-a-url`)
	_, err := parseToolDirectory(dir)
	if err == nil {
		t.Fatal("invalid repository was accepted")
	}
	for _, want := range []string{`tool "broken-tool"`, `title "Broken tool"`, `repository "not-a-url"`} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("invalid repository error %q does not name %q", err, want)
		}
	}
}

func TestC3BuildRejectsInvalidRepository(t *testing.T) {
	root := fixtureRoot(t, `---
title: Fixture post
slug: fixture-post
date: 2026-08-25
summary: A fixture post.
plain_summary: A fixture post.
tags:
  - method
draft: false
---

Fixture post.
`)
	if err := os.MkdirAll(filepath.Join(root, "static"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "static", "site.css"), []byte("body {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	articleCSS, err := os.ReadFile(filepath.Join(repositoryRoot(t), "static", "article_a1.css"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "static", "article_a1.css"), articleCSS, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := writeC3ToolAt(t, filepath.Join(root, "content", "tools", "broken-tool"), `repository: not-a-url`); err != nil {
		t.Fatal(err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(workingDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	err = Build(root, t.TempDir(), false)
	if err == nil {
		t.Fatal("build accepted invalid repository")
	}
	for _, want := range []string{`tool "broken-tool"`, `title "Broken tool"`, `repository "not-a-url"`} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("build error %q does not name %q", err, want)
		}
	}
}

func TestC3ValidRepositoryStillRendersLink(t *testing.T) {
	dir := writeC3Tool(t, `repository: https://example.com/published`)
	tool, err := parseToolDirectory(dir)
	if err != nil {
		t.Fatal(err)
	}
	output := renderComponent(t, toolPage(tool, nil, PageData{}, nil))
	if !bytes.Contains(output, []byte(`<a href="https://example.com/published" rel="noopener">Repository</a>`)) {
		t.Fatalf("valid repository link missing from rendered HTML: %s", output)
	}
}

func writeC3Tool(t *testing.T, repository string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "broken-tool")
	if _, err := writeC3ToolAt(t, dir, repository); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeC3ToolAt(t *testing.T, dir, repository string) (string, error) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	content := `---
title: Broken tool
slug: broken-tool
date: 2026-08-25
summary: A broken tool fixture.
language: Go
status: active
` + repository + `
install: from a local checkout
example: broken run
---
Why it exists.
`
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte(content), 0o644); err != nil {
		return "", err
	}
	return dir, nil
}
