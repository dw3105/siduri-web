package site

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func init() {
	RegisterContent("demo", func(data PageData, routes *RouteSet) {
		entries := make([]demoEntry, 0, len(data.Posts))
		for _, post := range data.Posts {
			entries = append(entries, demoEntry{Slug: post.Slug, Title: post.Title, Preview: data.Preview})
		}
		routes.Register(Route{
			Name: "demo-pages",
			Output: RouteOutput{Expand: func(_ PageData) []Output {
				outputs := make([]Output, 0, len(entries))
				for _, entry := range entries {
					entry := entry
					outputs = append(outputs, PageOutput(
						"demo/"+entry.Slug+"/index.html",
						func() interfaceComponent { return demoPage(entry, renderHeadFragments(data)) },
					))
				}
				return outputs
			}},
		})
		routes.Register(Route{
			Name: "demo-bytes",
			Output: RouteOutput{Expand: func(_ PageData) []Output {
				return []Output{ByteOutput("demo/tags.txt", []byte(strings.Join(data.Tags, ",")+"\n"))}
			}},
		})
	})
	RegisterHead(HeadFragment{
		Name:   "demo",
		Render: func(PageData) interfaceComponent { return demoHead() },
	})
}

func TestDemoLane(t *testing.T) {
	root := demoFixtureRoot(t)
	output := t.TempDir()
	if err := Build(root, output, false); err != nil {
		t.Fatal(err)
	}
	for _, slug := range []string{"first-demo", "second-demo"} {
		path := filepath.Join(output, "demo", slug, "index.html")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read fan-out page %s: %v", slug, err)
		}
		if !bytes.Contains(data, []byte(`<meta name="demo-head" content="registered">`)) {
			t.Fatalf("head fragment missing from %s", slug)
		}
	}
	tags, err := os.ReadFile(filepath.Join(output, "demo", "tags.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(tags) != "build-log,tool-release,dogfooding,method,outcome\n" {
		t.Fatalf("unexpected byte output %q", tags)
	}
}

func TestBuildTwiceWithContentRegistration(t *testing.T) {
	root := demoFixtureRoot(t)
	output := t.TempDir()
	for attempt := 1; attempt <= 2; attempt++ {
		if err := Build(root, output, false); err != nil {
			t.Fatalf("build %d: %v", attempt, err)
		}
	}
}

func TestDuplicateFanoutOutputPanics(t *testing.T) {
	set := newRouteSet()
	set.Register(Route{
		Name: "first",
		Output: RouteOutput{Expand: func(PageData) []Output {
			return []Output{ByteOutput("same.txt", []byte("first"))}
		}},
	})
	set.Register(Route{
		Name: "second",
		Output: RouteOutput{Expand: func(PageData) []Output {
			return []Output{ByteOutput("same.txt", []byte("second"))}
		}},
	})

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected duplicate fan-out output to panic")
		}
	}()
	outputsForBuild(set, PageData{})
}

func demoFixtureRoot(t *testing.T) string {
	t.Helper()
	root := fixtureRoot(t, `---
title: First demo
slug: first-demo
date: 2026-08-24
summary: A first demo page.
plain_summary: A first demo page.
tags:
  - method
draft: false
---

First demo body.`)
	if err := os.WriteFile(filepath.Join(root, "content", "posts", "second.md"), []byte(`---
title: Second demo
slug: second-demo
date: 2026-08-23
summary: A second demo page.
plain_summary: A second demo page.
tags:
  - outcome
draft: false
---

Second demo body.`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "static"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "static", "site.css"), []byte("body { color: black; }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
