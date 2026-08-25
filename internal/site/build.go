package site

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/a-h/templ"
)

func Build(root, output string, includeDrafts bool) error {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve project root: %w", err)
	}
	outputAbs, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("resolve output directory: %w", err)
	}
	if outputAbs == rootAbs {
		return fmt.Errorf("refusing to use project root as build output")
	}
	content, err := LoadContent(rootAbs)
	if err != nil {
		return err
	}
	posts := make([]Post, 0, len(content.Posts))
	for _, post := range content.Posts {
		if includeDrafts || !post.Draft {
			posts = append(posts, post)
		}
	}
	sort.Slice(posts, func(i, j int) bool {
		if posts[i].Date != posts[j].Date {
			return posts[i].Date > posts[j].Date
		}
		return posts[i].Slug < posts[j].Slug
	})
	if err := os.RemoveAll(outputAbs); err != nil {
		return fmt.Errorf("clear build output: %w", err)
	}
	if err := os.MkdirAll(outputAbs, 0o755); err != nil {
		return fmt.Errorf("create build output: %w", err)
	}
	css, err := os.ReadFile(filepath.Join(rootAbs, "static", "site.css"))
	if err != nil {
		return fmt.Errorf("read static stylesheet: %w", err)
	}
	if err := writeOutput(outputAbs, "site.css", css); err != nil {
		return err
	}
	data := PageData{Posts: posts, Preview: includeDrafts}
	for _, route := range registeredRoutes() {
		if err := writeComponent(outputAbs, route.Output, route.Render(data)); err != nil {
			return fmt.Errorf("render %s: %w", route.Name, err)
		}
	}
	for _, post := range posts {
		body := renderMarkdown(post.Body)
		path := filepath.Join("journal", post.Slug, "index.html")
		if err := writeComponent(outputAbs, path, ArticlePage(post, body, includeDrafts)); err != nil {
			return fmt.Errorf("render post %s: %w", post.Slug, err)
		}
	}
	return nil
}

type interfaceComponent = templ.Component

func writeComponent(root, relative string, component templ.Component) error {
	var rendered bytes.Buffer
	if err := component.Render(context.Background(), &rendered); err != nil {
		return err
	}
	return writeOutput(root, relative, append(rendered.Bytes(), '\n'))
}

func writeOutput(root, relative string, data []byte) error {
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output directory for %s: %w", relative, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relative, err)
	}
	return nil
}
