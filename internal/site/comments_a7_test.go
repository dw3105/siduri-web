package site

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestA7RepositoryCommentsRenderAsOneLevelThread(t *testing.T) {
	reference := time.Now().UTC()
	t.Setenv("SOURCE_DATE_EPOCH", strconv.FormatInt(reference.Unix(), 10))
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
	if !strings.Contains(output, `id="comment-01K3QZJ8X4YB7N2M9V0PQRSTUY"`) {
		t.Fatal("reader depth-2 comment is missing from rendered thread markup")
	}
	renderedCommentCount := strings.Count(output, `<div class="comment" id="comment-`)
	if renderedCommentCount == 0 {
		t.Fatal("rendered thread contains no comment elements")
	}
	commentCount := strconv.Itoa(renderedCommentCount)
	for _, want := range []string{
		`class="comment-thread"`,
		`data-comment-count="` + commentCount + `"`,
		`itemprop="commentCount" content="` + commentCount + `"`,
		`itemtype="https://schema.org/Comment"`,
		`Siduri reply`,
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
	thread := commentThreadForPost(Post{Slug: "not-in-repository", Date: commentFreezeReferenceDate().Format("2006-01-02")})
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

func TestA7RefusalDetailsUseOperatorAttributeAndReaderCopy(t *testing.T) {
	thread := commentThread{
		Refusals: []commentRefusal{
			{
				OperatorMessage: "Comments were not rendered: read comments for hello-siduri: permission denied",
				ReaderMessage:   "Comments are unavailable because the thread could not be loaded.",
			},
		},
	}
	output := string(renderComponent(t, commentsSection(Post{Slug: "hello-siduri"}, thread)))
	if !strings.Contains(output, `data-comment-refusal="Comments were not rendered: read comments for hello-siduri: permission denied"`) {
		t.Fatal("thread load detail is missing from the operator attribute")
	}
	if !strings.Contains(output, `>Comments are unavailable because the thread could not be loaded.</p>`) {
		t.Fatal("thread load refusal is missing its reader-facing copy")
	}
	if strings.Contains(output, `>Comments were not rendered:`) {
		t.Fatal("thread load detail leaked into reader-facing copy")
	}
}

func TestA7CommentFreezeUsesBuildReferenceDate(t *testing.T) {
	reference := time.Now().UTC()
	cutoff := reference.AddDate(-1, 0, 0)
	tests := []struct {
		date   string
		closed bool
	}{
		{date: cutoff.AddDate(0, 0, -1).Format("2006-01-02"), closed: true},
		{date: cutoff.Format("2006-01-02"), closed: false},
		{date: cutoff.AddDate(0, 0, 1).Format("2006-01-02"), closed: false},
	}
	for _, test := range tests {
		if got := commentsClosedAt(test.date, reference); got != test.closed {
			t.Errorf("commentsClosedAt(%q, reference) = %v, want %v", test.date, got, test.closed)
		}
	}
}

func TestA7CommentFreezeReferenceUsesSourceDateEpoch(t *testing.T) {
	t.Setenv("SOURCE_DATE_EPOCH", "0")
	want := time.Unix(0, 0).UTC()
	if got := commentFreezeReferenceDate(); !got.Equal(want) {
		t.Fatalf("commentFreezeReferenceDate() = %s, want %s", got, want)
	}
}

func TestA7CommentFreezeReferenceUsesWallClockWhenUnpinned(t *testing.T) {
	previous, hadPrevious := os.LookupEnv("SOURCE_DATE_EPOCH")
	if err := os.Unsetenv("SOURCE_DATE_EPOCH"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if hadPrevious {
			_ = os.Setenv("SOURCE_DATE_EPOCH", previous)
		} else {
			_ = os.Unsetenv("SOURCE_DATE_EPOCH")
		}
	})
	before := time.Now().UTC()
	got := commentFreezeReferenceDate()
	after := time.Now().UTC()
	if got.Before(before) || got.After(after) {
		t.Fatalf("commentFreezeReferenceDate() = %s, outside wall-clock interval %s–%s", got, before, after)
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
	xssFixture := `<script>alert("not executable")</script>`
	if got, want := renderCommentMarkdown(xssFixture), "<p>&lt;script&gt;alert(&#34;not executable&#34;)&lt;/script&gt;</p>\n"; got != want {
		t.Fatalf("renderCommentMarkdown(%q) = %q, want %q", xssFixture, got, want)
	}
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

func TestA7CommentArrangementMakesEveryRecordReachableOrRefused(t *testing.T) {
	tests := []struct {
		name               string
		records            []commentRecord
		wantRendered       []string
		wantRefused        []string
		wantRefusalDetails map[string]string
	}{
		{
			name: "missing parent chain",
			records: []commentRecord{
				testCommentRecord("B", "missing"),
				testCommentRecord("C", "B"),
			},
			wantRefused: []string{"B", "C"},
			wantRefusalDetails: map[string]string{
				"B": "its parent comment missing is missing.",
				"C": "its parent comment B was not rendered.",
			},
		},
		{
			name: "missing parent chain with another descendant",
			records: []commentRecord{
				testCommentRecord("B", "missing"),
				testCommentRecord("C", "B"),
				testCommentRecord("D", "C"),
			},
			wantRefused: []string{"B", "C", "D"},
			wantRefusalDetails: map[string]string{
				"B": "its parent comment missing is missing.",
				"C": "its parent comment B was not rendered.",
				"D": "its parent comment C was not rendered.",
			},
		},
		{
			name: "too deep after two allowed levels",
			records: []commentRecord{
				testCommentRecord("root", ""),
				testCommentRecord("A", "root"),
				testCommentRecord("B", "A"),
				testCommentRecord("C", "B"),
			},
			wantRendered: []string{"root", "A", "B"},
			wantRefused:  []string{"C"},
			wantRefusalDetails: map[string]string{
				"C": "replies can be attached only up to two levels deep.",
			},
		},
		{
			name: "one reply",
			records: []commentRecord{
				testCommentRecord("A", ""),
				testCommentRecord("B", "A"),
			},
			wantRendered: []string{"A", "B"},
		},
		{
			name: "siblings",
			records: []commentRecord{
				testCommentRecord("A", ""),
				testCommentRecord("B", "A"),
				testCommentRecord("B2", "A"),
			},
			wantRendered: []string{"A", "B", "B2"},
		},
		{
			name: "self parent",
			records: []commentRecord{
				testCommentRecord("A", "A"),
			},
			wantRefused: []string{"A"},
		},
		{
			name: "two node cycle",
			records: []commentRecord{
				testCommentRecord("A", "B"),
				testCommentRecord("B", "A"),
			},
			wantRefused: []string{"A", "B"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			comments, refusals := arrangeCommentRecords(test.records)
			rendered := testRenderedCommentCounts(comments)
			refused := testRefusedCommentCounts(refusals, test.records)

			for _, record := range test.records {
				renderedCount := rendered[record.ID]
				refusedCount := refused[record.ID]
				if !((renderedCount == 1 && refusedCount == 0) || (renderedCount == 0 && refusedCount == 1)) {
					t.Errorf("comment %s has rendered count %d and refusal count %d; want exactly one", record.ID, renderedCount, refusedCount)
				}
			}

			if got, want := sortedCommentIDs(rendered), sortedCommentIDsFromList(test.wantRendered); !slicesEqual(got, want) {
				t.Errorf("rendered comment IDs = %v, want %v", got, want)
			}
			if got, want := sortedCommentIDs(refused), sortedCommentIDsFromList(test.wantRefused); !slicesEqual(got, want) {
				t.Errorf("refused comment IDs = %v, want %v", got, want)
			}
			for id, detail := range test.wantRefusalDetails {
				if got := refusalDetailForComment(refusals, id); got != detail {
					t.Errorf("refusal detail for %s = %q, want %q", id, got, detail)
				}
			}
		})
	}
}

func testCommentRecord(id, parentID string) commentRecord {
	return commentRecord{
		ID:         id,
		AuthorName: "Reader " + id,
		AuthorRole: "reader",
		ParentID:   parentID,
		Body:       "Comment " + id,
	}
}

func testRenderedCommentCounts(comments []commentView) map[string]int {
	counts := make(map[string]int)
	var visit func([]commentView)
	visit = func(views []commentView) {
		for _, view := range views {
			counts[view.ID]++
			visit(view.Replies)
		}
	}
	visit(comments)
	return counts
}

func testRefusedCommentCounts(refusals []commentRefusal, records []commentRecord) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		prefix := "Comment " + record.ID + " was not rendered:"
		for _, refusal := range refusals {
			if strings.HasPrefix(refusal.OperatorMessage, prefix) {
				counts[record.ID]++
			}
		}
	}
	return counts
}

func sortedCommentIDs(counts map[string]int) []string {
	ids := make([]string, 0, len(counts))
	for id := range counts {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func sortedCommentIDsFromList(ids []string) []string {
	sorted := append([]string(nil), ids...)
	sort.Strings(sorted)
	return sorted
}

func slicesEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func refusalDetailForComment(refusals []commentRefusal, id string) string {
	prefix := "Comment " + id + " was not rendered: "
	for _, refusal := range refusals {
		if strings.HasPrefix(refusal.OperatorMessage, prefix) {
			return strings.TrimPrefix(refusal.OperatorMessage, prefix)
		}
	}
	return ""
}
