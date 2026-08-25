package site

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestA6ToolLoaderAndPage(t *testing.T) {
	tools, err := loadTools(filepath.Join("..", "..", "content", "tools"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tools) != 1 {
		t.Fatalf("expected one seeded tool, got %d", len(tools))
	}
	if tools[0].Slug != "gateslot" || tools[0].Status != "active" {
		t.Fatalf("unexpected seeded tool %#v", tools[0])
	}
	output := renderComponent(t, toolPage(tools[0], nil, PageData{}, nil))
	for _, want := range []string{"Gateslot", "Status: active", "Why it exists", "gateslot check"} {
		if !bytes.Contains(output, []byte(want)) {
			t.Fatalf("tool page missing %q", want)
		}
	}
	if bytes.Contains(output, []byte("Repository")) {
		t.Fatal("tool without a repository renders a repository label")
	}
}

func TestA6NoJavaScriptIndexListsEveryTool(t *testing.T) {
	tools := []Tool{
		{Title: "Recent", Slug: "recent", Date: "2026-08-25", Summary: "Recent tool", Language: "Go", Status: "active"},
		{Title: "Finished", Slug: "finished", Date: "2026-08-24", Summary: "Finished tool", Language: "Rust", Status: "abandoned"},
	}
	output := renderComponent(t, toolsIndexPage(tools, map[string][]Post{}, PageData{}, "", "", nil))
	if got := strings.Count(string(output), `class="tool-card"`); got != len(tools) {
		t.Fatalf("no-JavaScript index lists %d cards, want %d", got, len(tools))
	}
	if !bytes.Contains(output, []byte(`Status: abandoned`)) {
		t.Fatal("abandoned status is not visible on the index")
	}
	if bytes.Contains(output, []byte("<script")) {
		t.Fatal("tool index requires JavaScript")
	}
}

func TestA6AbandonedStatusIsVisibleOnToolPage(t *testing.T) {
	tool := Tool{
		Title: "Finished", Slug: "finished", Date: "2026-08-24", Summary: "Finished tool",
		Language: "Rust", Status: "abandoned", Repository: "https://example.invalid/finished",
		Install: "cargo install finished", Example: "finished run", Why: "The experiment is complete.",
	}
	output := renderComponent(t, toolPage(tool, nil, PageData{}, nil))
	if !bytes.Contains(output, []byte("Status: abandoned")) {
		t.Fatal("abandoned status is not visible on the tool page")
	}
}

func TestA6DanglingReferenceNamesPostAndSlug(t *testing.T) {
	_, err := postsByTool(
		[]Tool{{Slug: "gateslot"}},
		[]Post{{Slug: "uses-gateslot", Body: "This names tool:nope."}},
	)
	if err == nil || !strings.Contains(err.Error(), `post "uses-gateslot" references unknown tool "nope"`) {
		t.Fatalf("expected named dangling reference error, got %v", err)
	}
}

func TestA6BuildRejectsDanglingReference(t *testing.T) {
	root := fixtureRoot(t, `---
title: Broken tool reference
slug: broken-tool-reference
date: 2026-08-25
summary: A post with a dangling tool reference.
plain_summary: A post with a dangling tool reference.
tags:
  - method
draft: false
---

This post mentions tool:nope.`)
	if err := os.MkdirAll(filepath.Join(root, "static"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "static", "site.css"), []byte("body {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Build(root, t.TempDir(), false)
	if err == nil || !strings.Contains(err.Error(), `post "broken-tool-reference" references unknown tool "nope"`) {
		t.Fatalf("expected build error naming post and tool, got %v", err)
	}
}

func TestA6ReverseCrossReference(t *testing.T) {
	tools := []Tool{{
		Title: "Gateslot", Slug: "gateslot", Date: "2026-08-25", Summary: "An approval gate",
		Language: "Go", Status: "active", Repository: "https://example.invalid/gateslot",
		Install: "go install gateslot", Example: "gateslot check", Why: "It records the gate.",
	}}
	posts := []Post{{Title: "Using Gateslot", Slug: "using-gateslot", Body: "I used [Gateslot](/tools/gateslot/) in the build."}}
	mentioned, err := postsByTool(tools, posts)
	if err != nil {
		t.Fatal(err)
	}
	if len(mentioned["gateslot"]) != 1 || mentioned["gateslot"][0].Slug != "using-gateslot" {
		t.Fatalf("unexpected reverse references %#v", mentioned)
	}
	page := renderComponent(t, toolPage(tools[0], mentioned["gateslot"], PageData{}, nil))
	if !bytes.Contains(page, []byte(`/journal/using-gateslot/`)) {
		t.Fatal("tool page did not link the mentioning post")
	}
}

func TestA6TwoBuildsHaveIdenticalToolOutputs(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	first, second := t.TempDir(), t.TempDir()
	if err := Build(root, first, true); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, second, true); err != nil {
		t.Fatal(err)
	}
	firstFiles := readTree(t, first)
	secondFiles := readTree(t, second)
	if len(firstFiles) != len(secondFiles) {
		t.Fatalf("build file counts differ: %d and %d", len(firstFiles), len(secondFiles))
	}
	for path, firstData := range firstFiles {
		if !bytes.Equal(firstData, secondFiles[path]) {
			t.Fatalf("build output differs at %s", path)
		}
	}
}

func readTree(t *testing.T, root string) map[string][]byte {
	t.Helper()
	result := make(map[string][]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[relative] = data
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
