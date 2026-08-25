package site

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestA3FeedsAndMachineSurfaces(t *testing.T) {
	root := a3FixtureRoot(t)
	firstOutput := filepath.Join(t.TempDir(), "first")
	secondOutput := filepath.Join(t.TempDir(), "second")
	if err := Build(root, firstOutput, false); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, secondOutput, false); err != nil {
		t.Fatal(err)
	}

	rssData := a3ReadFile(t, firstOutput, "feed.xml")
	var rss struct {
		Channel struct {
			Items []struct {
				Title   string `xml:"title"`
				Encoded string `xml:"encoded"`
			} `xml:"item"`
		} `xml:"channel"`
	}
	if err := xml.Unmarshal(rssData, &rss); err != nil {
		t.Fatalf("feed.xml is not XML: %v", err)
	}
	if len(rss.Channel.Items) != 1 || rss.Channel.Items[0].Title != "A3 special characters" {
		t.Fatalf("unexpected RSS items: %+v", rss.Channel.Items)
	}
	if !strings.Contains(rss.Channel.Items[0].Encoded, "&lt;tag&gt;") || !strings.Contains(rss.Channel.Items[0].Encoded, "&amp; marker") {
		t.Fatalf("RSS content did not preserve escaped HTML: %q", rss.Channel.Items[0].Encoded)
	}

	var feed struct {
		Items []struct {
			URL         string `json:"url"`
			ContentHTML string `json:"content_html"`
		} `json:"items"`
	}
	if err := json.Unmarshal(a3ReadFile(t, firstOutput, "feed.json"), &feed); err != nil {
		t.Fatalf("feed.json is not JSON: %v", err)
	}
	if len(feed.Items) != 1 || feed.Items[0].URL != siteURL("/journal/a3-special/") {
		t.Fatalf("unexpected JSON Feed items: %+v", feed.Items)
	}
	if !strings.Contains(feed.Items[0].ContentHTML, "&lt;tag&gt;") || !strings.Contains(feed.Items[0].ContentHTML, "&amp; marker") {
		t.Fatalf("JSON Feed content did not preserve escaped HTML: %q", feed.Items[0].ContentHTML)
	}

	var sitemap struct {
		URLs []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}
	if err := xml.Unmarshal(a3ReadFile(t, firstOutput, "sitemap.xml"), &sitemap); err != nil {
		t.Fatalf("sitemap.xml is not XML: %v", err)
	}
	if !a3ContainsLoc(sitemap.URLs, siteURL("/journal/a3-special/")) {
		t.Fatalf("published post is missing from sitemap: %+v", sitemap.URLs)
	}

	for _, name := range []string{"feed.xml", "feed.json", "sitemap.xml", "llms.txt"} {
		if strings.Contains(string(a3ReadFile(t, firstOutput, name)), "a3-draft") {
			t.Fatalf("draft slug appears in %s", name)
		}
	}
	if !strings.Contains(string(a3ReadFile(t, firstOutput, "llms.txt")), "a3-special") {
		t.Fatal("published post is missing from llms.txt")
	}
	robots := string(a3ReadFile(t, firstOutput, "robots.txt"))
	if !strings.Contains(robots, "User-agent: *\nAllow: /\nSitemap: "+siteURL("/sitemap.xml")) {
		t.Fatalf("robots.txt has no sitemap declaration: %q", robots)
	}

	for _, name := range []string{"feed.xml", "feed.json"} {
		if first, second := a3ReadFile(t, firstOutput, name), a3ReadFile(t, secondOutput, name); !bytes.Equal(first, second) {
			t.Fatalf("%s changed between identical builds", name)
		}
	}
}

func TestA3PreviewStillExcludesDraftSurfaces(t *testing.T) {
	root := a3FixtureRoot(t)
	output := t.TempDir()
	if err := Build(root, output, true); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"feed.xml", "feed.json", "sitemap.xml", "llms.txt"} {
		if strings.Contains(string(a3ReadFile(t, output, name)), "a3-draft") {
			t.Fatalf("preview draft slug appears in %s", name)
		}
	}
}

func a3ContainsLoc(urls []struct {
	Loc string `xml:"loc"`
}, want string) bool {
	for _, url := range urls {
		if url.Loc == want {
			return true
		}
	}
	return false
}

func a3ReadFile(t *testing.T, root, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return data
}

func a3FixtureRoot(t *testing.T) string {
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
	published := `---
title: A3 special characters
slug: a3-special
date: 2026-08-24
summary: A summary for the A3 special-character fixture.
plain_summary: A post used to verify machine-readable surfaces.
tags:
  - method
draft: false
---

Literal <tag> and & marker.
`
	draft := `---
title: A3 draft
slug: a3-draft
date: 2026-08-25
summary: This draft must never enter a machine-readable surface.
plain_summary: This draft is intentionally excluded.
tags:
  - build-log
draft: true
---

Draft body.
`
	if err := os.WriteFile(filepath.Join(root, "content", "posts", "published.md"), []byte(published), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "content", "posts", "draft.md"), []byte(draft), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
