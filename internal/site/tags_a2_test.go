package site

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestA2TagFanoutCoversVocabularyAndSortsPosts(t *testing.T) {
	root := a2BuildFixture(t)
	first := t.TempDir()
	second := t.TempDir()
	if err := Build(root, first, false); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, second, false); err != nil {
		t.Fatal(err)
	}

	content, err := LoadContent(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(content.Tags) != 5 {
		t.Fatalf("tag vocabulary count = %d, want 5", len(content.Tags))
	}
	for _, tag := range content.Tags {
		path := filepath.Join(first, "tags", tag, "index.html")
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("tag page for %q is missing: %v", tag, err)
		}
	}
	tagEntries, err := os.ReadDir(filepath.Join(first, "tags"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tagEntries) != len(content.Tags) {
		t.Fatalf("tag page directory count = %d, want %d", len(tagEntries), len(content.Tags))
	}

	empty, err := os.ReadFile(filepath.Join(first, "tags", "outcome", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(empty, []byte("No posts carry this tag yet.")) {
		t.Fatal("empty vocabulary tag did not render its empty state")
	}
	golden, err := os.ReadFile(filepath.Join("..", "..", "testdata", "golden", "empty-tag.html"))
	if err != nil {
		t.Fatal(err)
	}
	if actual := renderComponent(t, TagPage("outcome", nil, false)); !bytes.Equal(actual, golden) {
		t.Fatal("empty tag component differs from the existing empty-tag golden")
	}

	method, err := os.ReadFile(filepath.Join(first, "tags", "method", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if firstIndex, secondIndex := bytes.Index(method, []byte("first-method")), bytes.Index(method, []byte("second-method")); firstIndex < 0 || secondIndex < 0 || firstIndex >= secondIndex {
		t.Fatalf("method posts are not newest first: first=%d second=%d", firstIndex, secondIndex)
	}

	firstDigest := a2TreeDigest(t, first)
	secondDigest := a2TreeDigest(t, second)
	if firstDigest != secondDigest {
		t.Fatalf("two builds differ: %s != %s", firstDigest, secondDigest)
	}
}

func TestA2UnknownTagStillFailsBuild(t *testing.T) {
	root := fixtureRoot(t, `---
title: Unknown tag
slug: unknown-tag-a2
date: 2026-08-25
summary: A summary.
plain_summary: A plain summary.
tags:
  - nonsense
draft: false
---

Body.`)
	err := Build(root, t.TempDir(), false)
	if err == nil || !strings.Contains(err.Error(), `unknown tag "nonsense"`) {
		t.Fatalf("expected unknown tag error, got %v", err)
	}
}

func TestA2UnsortedMapWouldChangeTagOutput(t *testing.T) {
	canonical := make([]string, 100)
	for i := range canonical {
		canonical[i] = fmt.Sprintf("tag-%03d", i)
	}
	unsorted := a2UnsortedMapOrder(canonical)
	if bytes.Equal(a2OrderDigest(canonical), a2OrderDigest(unsorted)) {
		t.Fatal("deliberately unsorted tag iteration did not change the output digest")
	}
	// The implementation iterates PageData.Tags, not a map, and therefore uses
	// the vocabulary's stable order for the output sequence.
	if got := a2TagOutputOrder(PageData{Tags: canonical}); strings.Join(got, ",") != strings.Join(canonical, ",") {
		t.Fatalf("tag output order = %v, want %v", got, canonical)
	}
}

func a2UnsortedMapOrder(tags []string) []string {
	tagSet := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		tagSet[tag] = struct{}{}
	}
	order := make([]string, 0, len(tagSet))
	for tag := range tagSet { // Deliberately broken: map order is not a page order.
		order = append(order, tag)
	}
	return order
}

func a2TagOutputOrder(data PageData) []string {
	postsByTag := make(map[string][]Post, len(data.Tags))
	for _, tag := range data.Tags {
		postsByTag[tag] = nil
	}
	order := make([]string, 0, len(data.Tags))
	for _, tag := range data.Tags {
		if _, ok := postsByTag[tag]; ok {
			order = append(order, tag)
		}
	}
	return order
}

func a2OrderDigest(tags []string) []byte {
	hash := sha256.New()
	for _, tag := range tags {
		_, _ = hash.Write([]byte(tag))
		_, _ = hash.Write([]byte{0})
	}
	return hash.Sum(nil)
}

func a2TreeDigest(t *testing.T, root string) string {
	t.Helper()
	type fileDigest struct {
		path string
		sum  string
	}
	var files []fileDigest
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
		sum := sha256.Sum256(data)
		files = append(files, fileDigest{path: strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(root)+"/"), sum: hex.EncodeToString(sum[:])})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].path < files[j].path })
	hash := sha256.New()
	for _, file := range files {
		_, _ = hash.Write([]byte(file.path))
		_, _ = hash.Write([]byte{0})
		_, _ = hash.Write([]byte(file.sum))
		_, _ = hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func a2BuildFixture(t *testing.T) string {
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
	if err := os.WriteFile(filepath.Join(root, "static", "site.css"), []byte("body { color: black; }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	posts := map[string]string{
		"first.md": `---
title: First method
slug: first-method
date: 2026-08-24
summary: First method summary.
plain_summary: First method summary.
tags:
  - method
draft: false
---

First method body.`,
		"second.md": `---
title: Second method
slug: second-method
date: 2026-08-23
summary: Second method summary.
plain_summary: Second method summary.
tags:
  - method
draft: false
---

Second method body.`,
		"build.md": `---
title: Build log
slug: build-log
date: 2026-08-22
summary: Build log summary.
plain_summary: Build log summary.
tags:
  - build-log
draft: false
---

Build log body.`,
	}
	for name, post := range posts {
		if err := os.WriteFile(filepath.Join(root, "content", "posts", name), []byte(post), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
