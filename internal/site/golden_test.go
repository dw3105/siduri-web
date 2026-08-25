package site

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoldenPages(t *testing.T) {
	long := Post{
		Title:        "A long build log",
		Slug:         "long-build-log",
		Date:         "2026-08-20",
		PlainSummary: "A long fixture exercises the article renderer.",
		Tags:         []string{"build-log"},
		Body:         strings.Repeat("This sentence keeps the representative article above the long-post threshold. ", 180),
	}
	short := Post{
		Title:        "A short note",
		Slug:         "short-note",
		Date:         "2026-08-21",
		PlainSummary: "A short fixture checks the compact article shape.",
		Tags:         []string{"method"},
		Body:         "A small page is still a page.",
	}
	fixtures := []struct {
		name      string
		component func() interfaceComponent
	}{
		{name: "long-post.html", component: func() interfaceComponent { return ArticlePage(long, renderMarkdown(long.Body), false) }},
		{name: "short-post.html", component: func() interfaceComponent { return ArticlePage(short, renderMarkdown(short.Body), false) }},
		{name: "empty-tag.html", component: func() interfaceComponent { return EmptyTagPage("outcome", false) }},
	}

	goldenDir := filepath.Join("..", "..", "testdata", "golden")
	for _, fixture := range fixtures {
		fixture := fixture
		t.Run(fixture.name, func(t *testing.T) {
			actual := renderComponent(t, fixture.component())
			path := filepath.Join(goldenDir, fixture.name)
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err := os.MkdirAll(goldenDir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, actual, 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			expected, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read golden file: %v; regenerate with UPDATE_GOLDEN=1", err)
			}
			if !bytes.Equal(expected, actual) {
				t.Fatalf("golden output differs; regenerate with UPDATE_GOLDEN=1")
			}
		})
	}
}

func renderComponent(t *testing.T, component interfaceComponent) []byte {
	t.Helper()
	var output bytes.Buffer
	if err := component.Render(context.Background(), &output); err != nil {
		t.Fatal(err)
	}
	return append(output.Bytes(), '\n')
}
