package site

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestA3HeadMetadataIsInsideHead(t *testing.T) {
	root := a3FixtureRoot(t)
	output := t.TempDir()
	if err := Build(root, output, false); err != nil {
		t.Fatal(err)
	}
	page := string(a3ReadFile(t, output, filepath.ToSlash(filepath.Join("journal", "a3-special", "index.html"))))
	headStart := strings.Index(page, "<head>")
	headEnd := strings.Index(page, "</head>")
	if headStart < 0 || headEnd < 0 || headStart > headEnd {
		t.Fatalf("page has no valid head: %q", page)
	}
	head := page[headStart:headEnd]
	for _, want := range []string{
		`rel="alternate" type="application/rss+xml"`,
		`rel="alternate" type="application/feed+json"`,
		`rel="canonical" href="https://siduri.ai/journal/a3-special/"`,
		`name="description" content="A summary for the A3 special-character fixture."`,
		`type="application/ld+json"`,
	} {
		if !strings.Contains(head, want) {
			t.Fatalf("head is missing %q: %s", want, head)
		}
	}
	if strings.Contains(page[:headStart], `rel="canonical"`) || strings.Contains(page[headEnd:], `rel="canonical"`) {
		t.Fatal("canonical link escaped the document head")
	}
	if strings.Contains(head, "a3-draft") {
		t.Fatal("draft appeared in head metadata")
	}

	start := strings.Index(head, `<script type="application/ld+json">`) + len(`<script type="application/ld+json">`)
	end := strings.Index(head[start:], "</script>")
	if start < len(`<script type="application/ld+json">`) || end < 0 {
		t.Fatal("JSON-LD script is not closed")
	}
	var graph struct {
		Graph []struct {
			Type string `json:"@type"`
		} `json:"@graph"`
	}
	if err := json.Unmarshal([]byte(head[start:start+end]), &graph); err != nil {
		t.Fatalf("JSON-LD is not JSON: %v", err)
	}
	if !a3HasSchemaType(graph.Graph, "Person") || !a3HasSchemaType(graph.Graph, "Article") {
		t.Fatalf("JSON-LD lacks Person and Article: %+v", graph.Graph)
	}
}

func a3HasSchemaType(entities []struct {
	Type string `json:"@type"`
}, want string) bool {
	for _, entity := range entities {
		if entity.Type == want {
			return true
		}
	}
	return false
}

func TestA3HeadMetadataHasDraftSafeFallback(t *testing.T) {
	data := renderComponent(t, a3Head(PageData{Posts: []Post{{
		Title:   "Draft only",
		Slug:    "draft-only",
		Summary: "Draft summary",
		Draft:   true,
	}}}))
	if strings.Contains(string(data), "Draft summary") || strings.Contains(string(data), "draft-only") {
		t.Fatalf("draft leaked into fallback metadata: %s", data)
	}
}
