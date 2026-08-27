package site

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestA5RequiredPagesAndOperatorLinkGraph(t *testing.T) {
	output := t.TempDir()
	if err := Build(repositoryRoot(t), output, false); err != nil {
		t.Fatal(err)
	}

	required := []string{
		"index.html",
		"about/index.html",
		"contact/index.html",
		"stack/index.html",
		"404.html",
		"impressum/index.html",
		"datenschutz/index.html",
	}
	for _, relative := range required {
		if _, err := os.Stat(filepath.Join(output, relative)); err != nil {
			t.Errorf("required A5 output %s: %v", relative, err)
		}
	}

	for _, relative := range htmlFiles(t, output) {
		data, err := os.ReadFile(filepath.Join(output, relative))
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(bytes.ToLower(data), []byte("<form")) {
			t.Errorf("form element found in %s", relative)
		}
	}

	for _, page := range htmlFiles(t, output) {
		data, err := os.ReadFile(filepath.Join(output, page))
		if err != nil {
			t.Fatal(err)
		}
		for _, href := range []string{`href="/impressum/"`, `href="/datenschutz/"`} {
			if bytes.Contains(data, []byte(href)) {
				t.Errorf("%s unexpectedly links %s", page, href)
			}
		}
	}

	about, err := os.ReadFile(filepath.Join(output, "about/index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(about, []byte("Siduri is the tavern keeper")) {
		t.Error("about page does not contain the Siduri principle story")
	}
	if bytes.Contains(about, []byte("/services/")) || bytes.Contains(about, []byte("€")) {
		t.Error("about page contains deferred sales copy")
	}
}

func TestA5NonPostContentStaysOutOfJournal(t *testing.T) {
	output := t.TempDir()
	if err := Build(repositoryRoot(t), output, false); err != nil {
		t.Fatal(err)
	}
	journal, err := os.ReadFile(filepath.Join(output, "journal/index.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"home", "about", "stack"} {
		if bytes.Contains(bytes.ToLower(journal), []byte(name+".md")) {
			t.Errorf("content/%s.md leaked into the journal", name)
		}
	}
	// The feed is A3's and must exist. Asserting its absence was green only
	// while A3 was unmerged -- a lane cannot assert the absence of another
	// lane's deliverable. What this test is named for is that non-post content
	// stays out of the surfaces posts appear on, so check that instead.
	feed, err := os.ReadFile(filepath.Join(output, "feed.xml"))
	if err != nil {
		t.Fatalf("feed.xml missing: %v", err)
	}
	for _, name := range []string{"home", "about", "stack"} {
		if bytes.Contains(bytes.ToLower(feed), []byte(name+".md")) {
			t.Errorf("content/%s.md leaked into the feed", name)
		}
	}
}

func TestA5PublishedPostGetsDeterministicOGCard(t *testing.T) {
	root := a5BuildFixture(t)
	first := t.TempDir()
	second := t.TempDir()
	if err := Build(root, first, false); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, second, false); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join("journal", "published-note", "og.png")
	firstBytes, err := os.ReadFile(filepath.Join(first, path))
	if err != nil {
		t.Fatalf("read first OG card: %v", err)
	}
	secondBytes, err := os.ReadFile(filepath.Join(second, path))
	if err != nil {
		t.Fatalf("read second OG card: %v", err)
	}
	if !bytes.Equal(firstBytes, secondBytes) {
		t.Fatal("OG card differs across identical builds")
	}
	if len(firstBytes) == 0 {
		t.Fatal("OG card is empty")
	}
}

func htmlFiles(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(relative))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)
	return files
}

func a5BuildFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "content", "posts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "static"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "content", "tags.yml"), []byte("tags:\n  - build-log\n  - tool-release\n  - dogfooding\n  - method\n  - outcome\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	post := `---
title: Published note
slug: published-note
date: 2026-08-25
summary: A published note for the A5 card test.
plain_summary: A published note for the A5 card test.
series: The build series
tags:
  - build-log
draft: false
---

This post is deliberately published in the fixture.
`
	if err := os.WriteFile(filepath.Join(root, "content", "posts", "published.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "static", "site.css"), []byte("body { color: black; }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(filepath.Join(root, "..", ".."))
}
