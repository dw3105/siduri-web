package site

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var a5HrefPattern = regexp.MustCompile(`href="([^"]+)"`)

func TestA5RequiredPagesAndLegalLinkGraph(t *testing.T) {
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

	pages := htmlFiles(t, output)
	for _, page := range pages {
		for _, target := range []string{"impressum/index.html", "datenschutz/index.html"} {
			if !a5ReachableWithin(page, target, output, 2) {
				t.Errorf("%s cannot reach %s in two clicks", page, target)
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

func a5ReachableWithin(start, target, root string, maxDepth int) bool {
	type node struct {
		page  string
		depth int
	}
	queue := []node{{page: start}}
	seen := map[string]bool{start: true}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.page == target {
			return true
		}
		if current.depth == maxDepth {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, current.page))
		if err != nil {
			continue
		}
		for _, match := range a5HrefPattern.FindAllSubmatch(data, -1) {
			next := a5NormalizeInternalLink(string(match[1]))
			if next == "" || seen[next] {
				continue
			}
			seen[next] = true
			queue = append(queue, node{page: next, depth: current.depth + 1})
		}
	}
	return false
}

func a5NormalizeInternalLink(href string) string {
	if href == "/" {
		return "index.html"
	}
	if !strings.HasPrefix(href, "/") || strings.HasPrefix(href, "//") {
		return ""
	}
	href = strings.TrimPrefix(href, "/")
	if strings.HasSuffix(href, "/") {
		href += "index.html"
	}
	return filepath.ToSlash(href)
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
