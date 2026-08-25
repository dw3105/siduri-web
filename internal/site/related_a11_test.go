package site

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestA11RelatedPostsRenderThreeBuildVisiblePosts(t *testing.T) {
	current := Post{
		Title:        "Current post",
		Slug:         "current-post",
		Date:         "2026-08-25",
		PlainSummary: "The current post.",
		Tags:         []string{"method"},
		Series:       "shared-series",
		Part:         "1",
		Body:         "The body.",
	}
	posts := []Post{
		current,
		{Title: "Second post", Slug: "second-post", Date: "2026-08-24", PlainSummary: "The second post.", Tags: []string{"method"}, Series: "shared-series", Part: "2", Body: "The body."},
		{Title: "Third post", Slug: "third-post", Date: "2026-08-23", PlainSummary: "The third post.", Tags: []string{"method"}, Series: "shared-series", Part: "3", Body: "The body."},
		{Title: "Fourth post", Slug: "fourth-post", Date: "2026-08-22", PlainSummary: "The fourth post.", Tags: []string{"method"}, Series: "shared-series", Part: "4", Body: "The body."},
	}

	output := renderA11Article(t, current, posts)
	section := a11RelatedSectionHTML(t, output)
	if got := strings.Count(section, `class="related-post"`); got != 3 {
		t.Fatalf("related post count = %d, want 3: %s", got, section)
	}
	for _, slug := range []string{"second-post", "third-post", "fourth-post"} {
		if !strings.Contains(section, `/journal/`+slug+`/`) {
			t.Fatalf("related section is missing %q: %s", slug, section)
		}
	}
	if strings.Contains(section, `/journal/`+current.Slug+`/`) {
		t.Fatalf("related section linked to the current post: %s", section)
	}
}

func TestA11EmptyRelatedPostsRenderNoSection(t *testing.T) {
	post := Post{
		Title:        "Only post",
		Slug:         "only-post",
		Date:         "2026-08-25",
		PlainSummary: "The only post.",
		Tags:         []string{"method"},
		Body:         "The body.",
	}

	section := string(renderComponent(t, relatedPostsSection(nil)))
	if strings.TrimSpace(section) != "" {
		t.Fatalf("empty related set rendered markup: %q", section)
	}
	output := renderA11Article(t, post, []Post{post})
	if strings.Contains(output, `class="related-posts"`) || strings.Contains(output, "Related posts") {
		t.Fatalf("one-post article rendered an empty related section: %s", output)
	}
}

func TestA11RelatedSectionIgnoresRegistrationOrder(t *testing.T) {
	post := Post{
		Title:        "Ordering post",
		Slug:         "ordering-post",
		Date:         "2026-08-25",
		PlainSummary: "The ordering post.",
		Tags:         []string{"method"},
		Series:       "ordering-series",
		Part:         "1",
		Body:         "The body.",
	}
	posts := []Post{
		post,
		{Title: "Ordering next", Slug: "ordering-next", Date: "2026-08-24", PlainSummary: "The next post.", Tags: []string{"method"}, Series: "ordering-series", Part: "2", Body: "The body."},
	}
	registrations := a11Sections(t, "a1-series-navigation", "a7-comments", "a11-related-posts")

	first := renderA11ArticleWithRegistrations(t, post, posts, registrations)
	reverse := append([]ArticleSection(nil), registrations...)
	for left, right := 0, len(reverse)-1; left < right; left, right = left+1, right-1 {
		reverse[left], reverse[right] = reverse[right], reverse[left]
	}
	second := renderA11ArticleWithRegistrations(t, post, posts, reverse)
	if !bytes.Equal([]byte(first), []byte(second)) {
		t.Fatalf("registration order changed article output\nfirst:  %s\nsecond: %s", first, second)
	}

	for _, marker := range []string{`class="series-navigation"`, `class="comment-thread"`, `class="related-posts"`} {
		if !strings.Contains(first, marker) {
			t.Fatalf("article section %q missing from output: %s", marker, first)
		}
	}
	if strings.Index(first, `class="series-navigation"`) > strings.Index(first, `class="related-posts"`) ||
		strings.Index(first, `class="related-posts"`) > strings.Index(first, `class="comment-thread"`) {
		t.Fatalf("after-body sections are not stably ordered: %s", first)
	}
}

func TestA11SinglePostBuildIsDeterministicAndIntentionallyHasNoRelatedMarkup(t *testing.T) {
	root := repositoryRoot(t)
	firstDir := t.TempDir()
	secondDir := t.TempDir()
	if err := Build(root, firstDir, false); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, secondDir, false); err != nil {
		t.Fatal(err)
	}

	firstFiles := a11BuildFiles(t, firstDir)
	secondFiles := a11BuildFiles(t, secondDir)
	if len(firstFiles) != len(secondFiles) {
		t.Fatalf("identical builds produced different file counts: %d and %d", len(firstFiles), len(secondFiles))
	}
	for path, firstData := range firstFiles {
		secondData, ok := secondFiles[path]
		if !ok || !bytes.Equal(firstData, secondData) {
			t.Fatalf("identical builds differ at %s", path)
		}
	}

	page := string(firstFiles[filepath.ToSlash(filepath.Join("journal", "hello-siduri", "index.html"))])
	if got := strings.Count(strings.ToLower(page), "related"); got != 0 {
		t.Fatalf("one-post build contains %d related matches; empty related sets should render nothing", got)
	}
}

func renderA11Article(t *testing.T, post Post, posts []Post) string {
	t.Helper()
	data := PageData{Posts: posts}
	return string(renderComponent(t, articlePage(post, renderArticleBody(post, data), data, nil, articleSectionsFor(post, data))))
}

func renderA11ArticleWithRegistrations(t *testing.T, post Post, posts []Post, registrations []ArticleSection) string {
	t.Helper()
	restore := isolateArticleSections()
	defer restore()
	for _, registration := range registrations {
		RegisterArticleSection(registration)
	}
	return renderA11Article(t, post, posts)
}

func a11RelatedSectionHTML(t *testing.T, output string) string {
	t.Helper()
	start := strings.Index(output, `<section class="related-posts"`)
	if start < 0 {
		t.Fatalf("related section is missing: %s", output)
	}
	end := strings.Index(output[start:], `</section>`)
	if end < 0 {
		t.Fatalf("related section is not closed: %s", output)
	}
	return output[start : start+end+len(`</section>`)]
}

func a11Sections(t *testing.T, names ...string) []ArticleSection {
	t.Helper()
	articleSectionMu.RLock()
	defer articleSectionMu.RUnlock()
	sections := make([]ArticleSection, 0, len(names))
	for _, name := range names {
		section, ok := articleSections[name]
		if !ok {
			t.Fatalf("article section %q is not registered", name)
		}
		sections = append(sections, section)
	}
	return sections
}

func a11BuildFiles(t *testing.T, root string) map[string][]byte {
	t.Helper()
	files := make(map[string][]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(relative)] = data
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
