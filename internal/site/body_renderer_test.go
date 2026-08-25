package site

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestArticleBodyRendererFallsBackToMarkdown(t *testing.T) {
	restore := isolateBodyRenderer()
	defer restore()
	post := Post{Body: "# Heading\n\nThe default body."}
	if got, want := renderArticleBody(post, PageData{}), renderMarkdown(post.Body); got != want {
		t.Fatalf("fallback body differs from renderMarkdown\ngot:  %q\nwant: %q", got, want)
	}
}

func TestSecondBodyRendererPanics(t *testing.T) {
	restore := isolateBodyRenderer()
	defer restore()
	RegisterBodyRenderer(BodyRenderer{
		Name:   "first-lane",
		Render: func(Post, PageData) string { return "first" },
	})

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected a second body renderer registration to panic")
		}
	}()
	RegisterBodyRenderer(BodyRenderer{
		Name:   "second-lane",
		Render: func(Post, PageData) string { return "second" },
	})
}

func TestRegisteredBodyRendererReplacesArticleBody(t *testing.T) {
	restore := isolateBodyRenderer()
	defer restore()
	RegisterBodyRenderer(BodyRenderer{
		Name: "marker-lane",
		Render: func(Post, PageData) string {
			return `<p data-body-renderer="marker">BODY_RENDERER_MARKER</p>`
		},
	})

	root := fixtureRoot(t, `---
title: Body renderer fixture
slug: body-renderer
date: 2026-08-25
summary: A body renderer fixture.
plain_summary: A body renderer fixture.
tags:
  - method
draft: false
---

Default body.`)
	if err := os.MkdirAll(filepath.Join(root, "static"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "static", "site.css"), []byte("body {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	output := t.TempDir()
	if err := Build(root, output, false); err != nil {
		t.Fatal(err)
	}
	page, err := os.ReadFile(filepath.Join(output, "journal", "body-renderer", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(page, []byte("BODY_RENDERER_MARKER")) {
		t.Fatal("registered body renderer marker missing from article page")
	}
	if bytes.Contains(page, []byte("<p>Default body.</p>")) {
		t.Fatal("default Markdown body was not replaced")
	}
}

func isolateBodyRenderer() func() {
	bodyRendererMu.Lock()
	saved := bodyRenderer
	bodyRenderer = nil
	bodyRendererMu.Unlock()
	return func() {
		bodyRendererMu.Lock()
		bodyRenderer = saved
		bodyRendererMu.Unlock()
	}
}
