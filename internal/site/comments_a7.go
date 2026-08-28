package site

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

// commentFreezeReferenceDate returns the UTC date used to evaluate FR-22.
// SOURCE_DATE_EPOCH pins this reference for reproducible builds; without it,
// the build uses the current UTC date so the freeze follows the real calendar.
func commentFreezeReferenceDate() time.Time {
	if value, ok := os.LookupEnv("SOURCE_DATE_EPOCH"); ok {
		seconds, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("site: invalid SOURCE_DATE_EPOCH %q: %v", value, err))
		}
		return time.Unix(seconds, 0).UTC()
	}
	return time.Now().UTC()
}

const commentLinkRel = "nofollow ugc noopener"

var (
	commentLinkPattern   = regexp.MustCompile(`(!)?\[([^\]]+)\]\(([^\s)]+)\)`)
	commentStrongPattern = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	commentUnderStrong   = regexp.MustCompile(`__([^_]+)__`)
	commentCodePattern   = regexp.MustCompile("`([^`]+)`")
	commentEmphasis      = regexp.MustCompile(`\*([^*]+)\*`)
	commentUnderEmphasis = regexp.MustCompile(`_([^_]+)_`)
	commentEmailPattern  = regexp.MustCompile(`[a-z0-9._%+-]+@[a-z0-9.-]+`)
	commentIDPattern     = regexp.MustCompile(`^[0-9A-Z]{26}$`)
)

var commentRecordFields = map[string]bool{
	"id":             true,
	"author_name":    true,
	"author_website": true,
	"author_role":    true,
	"email_hash":     true,
	"parent_id":      true,
	"created_at":     true,
	"approved_at":    true,
	"approved_by":    true,
}

type commentRecord struct {
	ID            string
	AuthorName    string
	AuthorWebsite string
	AuthorRole    string
	EmailHash     string
	ParentID      string
	CreatedAt     string
	ApprovedAt    string
	ApprovedBy    string
	Body          string
}

type commentView struct {
	ID            string
	AuthorName    string
	AuthorWebsite string
	AuthorRole    string
	CreatedAt     string
	BodyText      string
	BodyHTML      string
	Replies       []commentView
}

type commentThread struct {
	Comments   []commentView
	Count      int
	Closed     bool
	Refusals   []commentRefusal
	Structured string
}

type commentRefusal struct {
	OperatorMessage string
	ReaderMessage   string
}

func init() {
	RegisterArticleSection(ArticleSection{
		Name: "a7-comments",
		Slot: ArticleSectionAfterBody,
		Render: func(post Post, data PageData) templ.Component {
			// ArticlePage is also a small rendering helper used outside Build. In
			// that case, render when this slug has committed records or is frozen;
			// otherwise leave unrelated standalone fixtures unchanged.
			if !postBelongsToBuild(post, data) && !commentDirectoryExists(post.Slug) && !commentsClosed(post.Date) {
				return templ.Raw("")
			}
			return renderCommentSection(post)
		},
	})
	RegisterHead(HeadFragment{
		Name: "a7-comments-css",
		Render: func(PageData) templ.Component {
			return commentsStylesheet()
		},
	})
	RegisterContent("a7-comments-assets", func(data PageData, routes *RouteSet) {
		css, err := readCommentsCSS()
		if err != nil {
			panic(err)
		}
		routes.Register(Route{
			Name: "a7-comments-css",
			Output: RouteOutput{Expand: func(PageData) []Output {
				return []Output{ByteOutput("comments_a7.css", css)}
			}},
		})
	})
}

func postBelongsToBuild(post Post, data PageData) bool {
	for _, candidate := range data.Posts {
		if candidate.Slug == post.Slug {
			return true
		}
	}
	return false
}

func renderCommentSection(post Post) templ.Component {
	thread := commentThreadForPost(post)
	return commentsSection(post, thread)
}

func commentThreadForPost(post Post) commentThread {
	records, err := loadCommentRecords(post.Slug)
	thread := commentThread{Closed: commentsClosed(post.Date)}
	if err != nil {
		thread.Refusals = append(thread.Refusals, commentRefusal{
			OperatorMessage: "Comments were not rendered: " + err.Error(),
			ReaderMessage:   "Comments are unavailable because the thread could not be loaded.",
		})
		return thread
	}

	thread.Comments, thread.Refusals = arrangeCommentRecords(records)
	thread.Count = countCommentViews(thread.Comments)
	thread.Structured = commentStructuredData(post, thread.Comments)
	return thread
}

func countCommentViews(comments []commentView) int {
	count := 0
	for _, comment := range comments {
		count++
		count += countCommentViews(comment.Replies)
	}
	return count
}

func commentsClosed(postDate string) bool {
	return commentsClosedAt(postDate, commentFreezeReferenceDate())
}

func commentsClosedAt(postDate string, reference time.Time) bool {
	date, err := time.Parse("2006-01-02", postDate)
	if err != nil {
		return false
	}
	reference = reference.UTC()
	referenceDate := time.Date(reference.Year(), reference.Month(), reference.Day(), 0, 0, 0, 0, time.UTC)
	return date.Before(referenceDate.AddDate(-1, 0, 0))
}

func arrangeCommentRecords(records []commentRecord) ([]commentView, []commentRefusal) {
	byID := make(map[string]commentRecord, len(records))
	for _, record := range records {
		byID[record.ID] = record
	}

	views := make(map[string]commentView, len(records))
	for _, record := range records {
		views[record.ID] = commentView{
			ID:            record.ID,
			AuthorName:    record.AuthorName,
			AuthorWebsite: safeCommentWebsite(record.AuthorWebsite),
			AuthorRole:    record.AuthorRole,
			CreatedAt:     record.CreatedAt,
			BodyText:      record.Body,
			BodyHTML:      renderCommentMarkdown(record.Body),
		}
	}

	rootIDs := make([]string, 0, len(records))
	children := make(map[string][]string, len(records))
	refusals := make([]commentRefusal, 0)
	for _, record := range records {
		if record.ParentID == "" {
			if record.AuthorRole != "reader" {
				refusals = append(refusals, commentRefusalForComment(record.ID, "a site reply needs a parent comment.", "This comment was not shown because a site reply needs a parent comment."))
				continue
			}
			rootIDs = append(rootIDs, record.ID)
			continue
		}
		parent, exists := byID[record.ParentID]
		if !exists {
			refusals = append(refusals, commentRefusalForComment(record.ID, fmt.Sprintf("its parent comment %s is missing.", record.ParentID), "This comment was not shown because its parent comment is missing."))
			continue
		}
		if parent.ParentID != "" {
			grandparent, exists := byID[parent.ParentID]
			if exists && grandparent.ParentID != "" {
				refusals = append(refusals, commentRefusalForComment(record.ID, "replies can be attached only up to two levels deep.", "This comment was not shown because replies can be attached only up to two levels deep."))
				continue
			}
		}
		children[parent.ID] = append(children[parent.ID], record.ID)
	}
	sort.Strings(rootIDs)
	comments := make([]commentView, 0, len(rootIDs))
	var assemble func(string) commentView
	assemble = func(id string) commentView {
		view := views[id]
		replyIDs := children[id]
		sort.Strings(replyIDs)
		for _, replyID := range replyIDs {
			view.Replies = append(view.Replies, assemble(replyID))
		}
		return view
	}
	for _, id := range rootIDs {
		comments = append(comments, assemble(id))
	}
	sort.Slice(refusals, func(i, j int) bool {
		return refusals[i].OperatorMessage < refusals[j].OperatorMessage
	})
	return comments, refusals
}

func commentRefusalForComment(id, detail, readerMessage string) commentRefusal {
	return commentRefusal{
		OperatorMessage: fmt.Sprintf("Comment %s was not rendered: %s", id, detail),
		ReaderMessage:   readerMessage,
	}
}

func loadCommentRecords(slug string) ([]commentRecord, error) {
	dir := filepath.Join(commentsProjectRoot(), "content", "comments", slug)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read comments for %s: %w", slug, err)
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			return nil, fmt.Errorf("comment directory contains nested directory %q", entry.Name())
		}
		if filepath.Ext(entry.Name()) != ".md" {
			return nil, fmt.Errorf("comment directory contains non-Markdown file %q", entry.Name())
		}
		paths = append(paths, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(paths)

	records := make([]commentRecord, 0, len(paths))
	for _, path := range paths {
		record, err := parseCommentRecord(path)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func commentDirectoryExists(slug string) bool {
	path := filepath.Join(commentsProjectRoot(), "content", "comments", slug)
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func parseCommentRecord(path string) (commentRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return commentRecord{}, fmt.Errorf("read comment %s: %w", path, err)
	}
	if commentEmailPattern.MatchString(bytesToLower(data)) {
		return commentRecord{}, fmt.Errorf("%s: raw email address is forbidden in a comment record", path)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return commentRecord{}, fmt.Errorf("%s: frontmatter must start with ---", path)
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return commentRecord{}, fmt.Errorf("%s: frontmatter has no closing ---", path)
	}
	fields, err := parseCommentFrontmatter(lines[1:end], path)
	if err != nil {
		return commentRecord{}, err
	}
	for field := range commentRecordFields {
		if _, ok := fields[field]; !ok {
			return commentRecord{}, fmt.Errorf("%s: missing required frontmatter field %q", path, field)
		}
	}
	if len(fields) != len(commentRecordFields) {
		return commentRecord{}, fmt.Errorf("%s: comment frontmatter must match the repository record fields exactly", path)
	}

	record := commentRecord{
		ID:            fields["id"],
		AuthorName:    fields["author_name"],
		AuthorWebsite: fields["author_website"],
		AuthorRole:    fields["author_role"],
		EmailHash:     fields["email_hash"],
		ParentID:      emptyNull(fields["parent_id"]),
		CreatedAt:     fields["created_at"],
		ApprovedAt:    fields["approved_at"],
		ApprovedBy:    fields["approved_by"],
		Body:          strings.TrimSpace(strings.Join(lines[end+1:], "\n")),
	}
	if record.ID == "" || record.AuthorName == "" || record.AuthorRole == "" || record.EmailHash == "" || record.CreatedAt == "" || record.ApprovedAt == "" || record.ApprovedBy == "" {
		return commentRecord{}, fmt.Errorf("%s: required comment fields must not be empty", path)
	}
	if !commentIDPattern.MatchString(record.ID) {
		return commentRecord{}, fmt.Errorf("%s: id must be a 26-character ULID", path)
	}
	if record.AuthorRole != "reader" && record.AuthorRole != "site" {
		return commentRecord{}, fmt.Errorf("%s: author_role must be reader or site", path)
	}
	if !strings.HasPrefix(record.EmailHash, "sha256:") {
		return commentRecord{}, fmt.Errorf("%s: email_hash must use the sha256: prefix", path)
	}
	if _, err := time.Parse(time.RFC3339, record.CreatedAt); err != nil {
		return commentRecord{}, fmt.Errorf("%s: created_at must be RFC3339: %w", path, err)
	}
	if _, err := time.Parse(time.RFC3339, record.ApprovedAt); err != nil {
		return commentRecord{}, fmt.Errorf("%s: approved_at must be RFC3339: %w", path, err)
	}
	wantID := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if record.ID != wantID {
		return commentRecord{}, fmt.Errorf("%s: id %q must match the filename", path, record.ID)
	}
	return record, nil
}

func bytesToLower(data []byte) string {
	return strings.ToLower(string(data))
}

func parseCommentFrontmatter(lines []string, path string) (map[string]string, error) {
	fields := make(map[string]string, len(commentRecordFields))
	for _, raw := range lines {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		colon := strings.IndexByte(trimmed, ':')
		if colon <= 0 {
			return nil, fmt.Errorf("%s: malformed frontmatter line %q", path, trimmed)
		}
		key := strings.TrimSpace(trimmed[:colon])
		if !commentRecordFields[key] {
			return nil, fmt.Errorf("%s: unknown comment frontmatter field %q", path, key)
		}
		if _, exists := fields[key]; exists {
			return nil, fmt.Errorf("%s: duplicate comment frontmatter field %q", path, key)
		}
		fields[key] = unquote(strings.TrimSpace(trimmed[colon+1:]))
	}
	return fields, nil
}

func emptyNull(value string) string {
	if value == "null" {
		return ""
	}
	return value
}

func safeCommentWebsite(website string) string {
	if website == "" || commentSafeURL(website) {
		return website
	}
	return ""
}

func commentSafeURL(url string) bool {
	return (strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "//")) || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")
}

func renderCommentMarkdown(source string) string {
	lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
	var out strings.Builder
	paragraph := make([]string, 0)
	inCode := false
	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(renderCommentInline(strings.Join(paragraph, "\n")))
		out.WriteString("</p>\n")
		paragraph = paragraph[:0]
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCode {
				out.WriteString("</code></pre>\n")
				inCode = false
				continue
			}
			flushParagraph()
			language := strings.TrimSpace(strings.TrimPrefix(line, "```"))
			out.WriteString("<pre><code")
			if language != "" {
				out.WriteString(` class="language-`)
				out.WriteString(html.EscapeString(language))
				out.WriteByte('"')
			}
			out.WriteByte('>')
			inCode = true
			continue
		}
		if inCode {
			out.WriteString(html.EscapeString(line))
			out.WriteByte('\n')
			continue
		}
		if strings.TrimSpace(line) == "" {
			flushParagraph()
			continue
		}
		paragraph = append(paragraph, line)
	}
	if inCode {
		out.WriteString("</code></pre>\n")
	} else {
		flushParagraph()
	}
	return out.String()
}

func renderCommentInline(source string) string {
	const marker = "\x00A7COMMENT-LINK-%d\x00"
	links := make([]string, 0)
	protected := commentLinkPattern.ReplaceAllStringFunc(source, func(match string) string {
		parts := commentLinkPattern.FindStringSubmatch(match)
		if parts[1] == "!" || !commentSafeURL(parts[3]) {
			return html.EscapeString(match)
		}
		links = append(links, `<a href="`+html.EscapeString(parts[3])+`" rel="`+commentLinkRel+`">`+renderCommentInline(parts[2])+`</a>`)
		return fmt.Sprintf(marker, len(links)-1)
	})
	escaped := html.EscapeString(protected)
	escaped = commentStrongPattern.ReplaceAllString(escaped, `<strong>$1</strong>`)
	escaped = commentUnderStrong.ReplaceAllString(escaped, `<strong>$1</strong>`)
	escaped = commentCodePattern.ReplaceAllString(escaped, `<code>$1</code>`)
	escaped = commentEmphasis.ReplaceAllString(escaped, `<em>$1</em>`)
	escaped = commentUnderEmphasis.ReplaceAllString(escaped, `<em>$1</em>`)
	for i, link := range links {
		escaped = strings.ReplaceAll(escaped, html.EscapeString(fmt.Sprintf(marker, i)), link)
	}
	return strings.ReplaceAll(escaped, "\n", "<br />\n")
}

type commentStructuredItem struct {
	Type          string             `json:"@type"`
	ID            string             `json:"@id"`
	Author        commentSchemaValue `json:"author"`
	DatePublished string             `json:"datePublished"`
	Text          string             `json:"text"`
	IsPartOf      commentSchemaValue `json:"isPartOf"`
}

type commentSchemaValue struct {
	Type string `json:"@type"`
	ID   string `json:"@id,omitempty"`
	Name string `json:"name,omitempty"`
}

type commentStructuredGraph struct {
	Context string                  `json:"@context"`
	Graph   []commentStructuredItem `json:"@graph"`
}

func commentStructuredData(post Post, comments []commentView) string {
	if len(comments) == 0 {
		return ""
	}
	items := make([]commentStructuredItem, 0, countCommentViews(comments))
	var appendItem func(commentView)
	appendItem = func(comment commentView) {
		items = append(items, commentStructuredItem{
			Type: "Comment",
			ID:   siteURL("/journal/" + post.Slug + "/#comment-" + comment.ID),
			Author: commentSchemaValue{
				Type: "Person",
				Name: comment.AuthorName,
			},
			DatePublished: comment.CreatedAt,
			Text:          comment.BodyText,
			IsPartOf: commentSchemaValue{
				Type: "Article",
				ID:   canonicalURL(post),
			},
		})
		for _, reply := range comment.Replies {
			appendItem(reply)
		}
	}
	for _, comment := range comments {
		appendItem(comment)
	}
	data, err := json.Marshal(commentStructuredGraph{Context: "https://schema.org", Graph: items})
	if err != nil {
		panic("site: encode comment JSON-LD")
	}
	return `<script type="application/ld+json">` + string(data) + `</script>`
}

func commentsProjectRoot() string {
	working, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if fileExists(filepath.Join(working, "content", "tags.yml")) {
			return working
		}
		parent := filepath.Dir(working)
		if parent == working {
			return "."
		}
		working = parent
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func readCommentsCSS() ([]byte, error) {
	paths := []string{filepath.Join(commentsProjectRoot(), "static", "comments_a7.css")}
	if _, source, _, ok := runtime.Caller(0); ok {
		paths = append(paths, filepath.Join(filepath.Dir(source), "..", "..", "static", "comments_a7.css"))
	}
	var lastErr error
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("read comments stylesheet: %w", lastErr)
}
