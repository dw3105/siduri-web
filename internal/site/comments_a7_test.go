package site

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestA7RepositoryCommentsRenderAsOneLevelThread(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	root := repositoryRoot(t)
	if err := Build(root, first, true); err != nil {
		t.Fatal(err)
	}
	if err := Build(root, second, true); err != nil {
		t.Fatal(err)
	}

	pagePath := filepath.Join("journal", "hello-siduri", "index.html")
	page, err := os.ReadFile(filepath.Join(first, pagePath))
	if err != nil {
		t.Fatal(err)
	}
	secondPage, err := os.ReadFile(filepath.Join(second, pagePath))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(page, secondPage) {
		t.Fatal("comment page changed between identical builds")
	}
	firstFiles := a7BuildFiles(t, first)
	secondFiles := a7BuildFiles(t, second)
	if len(firstFiles) != len(secondFiles) {
		t.Fatalf("identical builds produced different file counts: %d and %d", len(firstFiles), len(secondFiles))
	}
	for path, firstData := range firstFiles {
		if secondData, ok := secondFiles[path]; !ok || !bytes.Equal(firstData, secondData) {
			t.Fatalf("identical builds differ at %s", path)
		}
	}

	output := string(page)
	for _, want := range []string{
		`class="comment-thread"`,
		`data-comment-count="2"`,
		`itemprop="commentCount" content="2"`,
		`itemtype="https://schema.org/Comment"`,
		`Siduri reply`,
		`Comment 01K3QZJ8X4YB7N2M9V0PQRSTUX was not rendered: replies can be attached only to top-level comments.`,
		`&lt;script&gt;alert(&#34;not executable&#34;)&lt;/script&gt;`,
		`rel="nofollow ugc noopener"`,
		`"@type":"Comment"`,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("A7 page is missing %q", want)
		}
	}
	if strings.Contains(output, `<script>alert`) {
		t.Fatal("comment body emitted executable HTML")
	}
	if got := strings.Count(output, `rel="nofollow ugc noopener"`); got < 4 {
		t.Fatalf("expected the three body links and author website to carry rel, got %d", got)
	}
	if strings.Contains(output, "This reply to a reply must be refused") {
		t.Fatal("reply to a reply was silently rendered")
	}
	if _, err := os.Stat(filepath.Join(first, "comments_a7.css")); err != nil {
		t.Fatalf("comments stylesheet was not emitted: %v", err)
	}
}

func a7BuildFiles(t *testing.T, root string) map[string][]byte {
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

func TestA7ZeroCommentsHasNoEmptyStateNoise(t *testing.T) {
	thread := commentThreadForPost(Post{Slug: "not-in-repository", Date: commentFreezeReferenceDate})
	if thread.Count != 0 || len(thread.Comments) != 0 {
		t.Fatalf("unexpected comments in empty thread: %+v", thread)
	}
	output := string(renderComponent(t, commentsSection(Post{Slug: "not-in-repository"}, thread)))
	if !strings.Contains(output, `class="comment-thread"`) {
		t.Fatal("empty thread region is missing")
	}
	if strings.Contains(strings.ToLower(output), "no comments") {
		t.Fatalf("empty-state noise leaked into empty thread: %s", output)
	}
}

func TestA7CommentFreezeUsesFixedReferenceDate(t *testing.T) {
	tests := []struct {
		date   string
		closed bool
	}{
		{date: "2025-08-24", closed: true},
		{date: "2025-08-25", closed: false},
		{date: "2026-08-25", closed: false},
	}
	for _, test := range tests {
		if got := commentsClosed(test.date); got != test.closed {
			t.Errorf("commentsClosed(%q) = %v, want %v", test.date, got, test.closed)
		}
	}
}

func TestA7FrozenThreadShowsOneLineAndNeverShowsAForm(t *testing.T) {
	thread := commentThread{
		Comments: []commentView{{
			ID:         "01K3QZJ8X4YB7N2M9V0PQRSTUV",
			AuthorName: "Reader",
			CreatedAt:  "2025-01-01T00:00:00Z",
			BodyHTML:   "<p>Existing comment.</p>\n",
		}},
		Count:  1,
		Closed: true,
	}
	output := string(renderComponent(t, commentsSection(Post{Slug: "old-post"}, thread)))
	if !strings.Contains(output, "Comments are closed on this post because it is more than 12 months old.") {
		t.Fatal("frozen thread is missing its closed explanation")
	}
	if !strings.Contains(output, "Existing comment.") {
		t.Fatal("frozen thread dropped an existing comment")
	}
	if strings.Contains(output, "<form") {
		t.Fatal("frozen thread contains a comment form")
	}
}

func TestA7RestrictedMarkdownEscapesHTMLAndSupportsAllowedSubset(t *testing.T) {
	output := renderCommentMarkdown("one\nline\n\n**bold** *italic* `code` [link](https://example.com) ![image](https://example.com/image.png)\n\n# not a heading\n\n```go\n<script>\n```")
	for _, want := range []string{
		"<p>one<br />\nline</p>",
		"<strong>bold</strong>",
		"<em>italic</em>",
		"<code>code</code>",
		`<a href="https://example.com" rel="nofollow ugc noopener">link</a>`,
		"&lt;script&gt;",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("restricted renderer missing %q in %s", want, output)
		}
	}
	if strings.Contains(output, "<h1>") || strings.Contains(output, "<table>") || strings.Contains(output, "<img") || strings.Contains(output, "<script>") {
		t.Fatalf("restricted renderer emitted a forbidden element: %s", output)
	}
	if strings.Contains(output, `href="https://example.com/image.png"`) {
		t.Fatalf("image syntax was turned into a link: %s", output)
	}
}
