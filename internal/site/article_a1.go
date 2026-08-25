package site

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/a-h/templ"
)

const a1TOCWordThreshold = 1500

// Reading time uses the shared fixed divisor of 200 words per minute, keeping
// the visible estimate stable across builds.

// a1Heading is derived from the same Markdown source that produces the
// article body. IDs are therefore stable without asking an author to maintain
// a second list of anchors.
type a1Heading struct {
	ID    string
	Level int
	Text  string
}

type a1SeriesNav struct {
	Series       string
	Position     string
	Previous     *Post
	Next         *Post
	PreviousText string
	NextText     string
}

func init() {
	RegisterBodyRenderer(BodyRenderer{
		Name:   "a1-article",
		Render: func(post Post, _ PageData) string { return renderA1Markdown(post.Body) },
	})
	RegisterArticleSection(ArticleSection{
		Name: "a1-table-of-contents",
		Slot: ArticleSectionBodyAside,
		Render: func(post Post, _ PageData) templ.Component {
			_, headings := renderA1MarkdownWithHeadings(post.Body)
			if len(strings.Fields(post.Body)) <= a1TOCWordThreshold {
				headings = nil
			}
			return a1TableOfContents(headings)
		},
	})
	RegisterArticleSection(ArticleSection{
		Name: "a1-series-navigation",
		Slot: ArticleSectionAfterBody,
		Render: func(post Post, data PageData) templ.Component {
			return a1SeriesNavigation(seriesNavigationFor(post, data.Posts))
		},
	})
	RegisterArticleSection(ArticleSection{
		Name: "a1-contextual-slot",
		Slot: ArticleSectionFooter,
		Render: func(Post, PageData) templ.Component {
			// FR-9 reserves this slot for a future contextual line. It is
			// intentionally empty until the phase-3 offer exists.
			return a1ContextualSlot()
		},
	})
	RegisterHead(HeadFragment{
		Name: "a1-article-assets",
		Render: func(PageData) templ.Component {
			return templ.Raw(`<link rel="stylesheet" href="/article_a1.css"/>`)
		},
	})
	RegisterContent("a1-article-assets", registerA1Assets)
}

func registerA1Assets(_ PageData, routes *RouteSet) {
	css, err := os.ReadFile(filepath.Join(a1ProjectRoot(), "static", "article_a1.css"))
	if err != nil {
		panic(fmt.Errorf("a1 article assets: %w", err))
	}
	routes.Register(Route{
		Name: "a1-article-stylesheet",
		Output: RouteOutput{Expand: func(PageData) []Output {
			return []Output{ByteOutput("article_a1.css", css)}
		}},
	})
}

// a1ProjectRoot mirrors the root discovery needed by the public registration
// seam. Build's root is deliberately not added to PageData, so this keeps the
// asset lookup local to the lane and works when go test runs in internal/site.
func a1ProjectRoot() string {
	working, err := os.Getwd()
	if err != nil {
		return "."
	}
	for directory := working; ; directory = filepath.Dir(directory) {
		if _, err := os.Stat(filepath.Join(directory, "static", "article_a1.css")); err == nil {
			return directory
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			break
		}
	}
	return working
}

func renderA1Markdown(source string) string {
	rendered, _ := renderA1MarkdownWithHeadings(source)
	return rendered
}

func renderA1MarkdownWithHeadings(source string) (string, []a1Heading) {
	lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
	var out strings.Builder
	paragraph := make([]string, 0)
	headings := make([]a1Heading, 0)
	usedIDs := make(map[string]int)
	inCode := false
	var fence a1Fence
	var codeLines []string

	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(renderInline(strings.Join(paragraph, "\n")))
		out.WriteString("</p>\n")
		paragraph = paragraph[:0]
	}
	flushCode := func() {
		out.WriteString(a1CodeBlockHTML(fence, strings.Join(codeLines, "\n")))
		codeLines = nil
		fence = a1Fence{}
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCode {
				flushCode()
				inCode = false
			} else {
				flushParagraph()
				fence = parseA1Fence(strings.TrimSpace(strings.TrimPrefix(line, "```")))
				inCode = true
			}
			continue
		}
		if inCode {
			codeLines = append(codeLines, line)
			continue
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			flushParagraph()
			continue
		}
		if level, text, ok := a1MarkdownHeading(trimmed); ok {
			flushParagraph()
			id := a1HeadingID(text, usedIDs)
			headings = append(headings, a1Heading{ID: id, Level: level, Text: a1PlainHeadingText(text)})
			fmt.Fprintf(&out, `<h%d id="%s">%s</h%d>`+"\n", level, html.EscapeString(id), renderInline(text), level)
			continue
		}
		if strings.HasPrefix(trimmed, "- ") {
			flushParagraph()
			out.WriteString("<ul>\n<li>")
			out.WriteString(renderInline(strings.TrimPrefix(trimmed, "- ")))
			out.WriteString("</li>\n</ul>\n")
			continue
		}
		paragraph = append(paragraph, line)
	}
	if inCode {
		flushCode()
	} else {
		flushParagraph()
	}
	return out.String(), headings
}

type a1Fence struct {
	language string
	filename string
}

func parseA1Fence(info string) a1Fence {
	info = strings.TrimSpace(info)
	if info == "" {
		return a1Fence{}
	}
	fields := strings.Fields(info)
	fence := a1Fence{language: strings.ToLower(fields[0])}
	rest := strings.TrimSpace(strings.TrimPrefix(info, fields[0]))
	for _, key := range []string{"filename", "file", "path", "title"} {
		marker := key + "="
		if index := strings.Index(strings.ToLower(rest), marker); index >= 0 {
			value := strings.TrimSpace(rest[index+len(marker):])
			if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
				if unquoted, err := strconv.Unquote(value); err == nil {
					value = unquoted
				}
			}
			fence.filename = strings.Trim(value, "'\"")
			break
		}
	}
	if fence.filename == "" && rest != "" && !strings.Contains(rest, "=") {
		fence.filename = strings.Trim(rest, "'\"")
	}
	return fence
}

func a1MarkdownHeading(line string) (int, string, bool) {
	switch {
	case strings.HasPrefix(line, "### "):
		return 3, strings.TrimSpace(strings.TrimPrefix(line, "### ")), true
	case strings.HasPrefix(line, "## "):
		return 2, strings.TrimSpace(strings.TrimPrefix(line, "## ")), true
	case strings.HasPrefix(line, "# "):
		return 2, strings.TrimSpace(strings.TrimPrefix(line, "# ")), true
	default:
		return 0, "", false
	}
}

func a1HeadingID(text string, used map[string]int) string {
	var id strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			id.WriteRune(r)
			lastDash = false
			continue
		}
		if id.Len() > 0 && !lastDash {
			id.WriteByte('-')
			lastDash = true
		}
	}
	base := strings.Trim(id.String(), "-")
	if base == "" {
		base = "section"
	}
	count := used[base]
	used[base] = count + 1
	if count > 0 {
		return base + "-" + strconv.Itoa(count+1)
	}
	return base
}

func a1PlainHeadingText(text string) string {
	text = strings.ReplaceAll(text, "**", "")
	text = strings.ReplaceAll(text, "__", "")
	text = strings.ReplaceAll(text, "`", "")
	return strings.TrimSpace(text)
}

func a1CodeBlockHTML(fence a1Fence, source string) string {
	known := a1KnownLanguage(fence.language)
	code := html.EscapeString(source + "\n")
	if known {
		code = a1HighlightCode(fence.language, source+"\n")
	}
	class := ""
	if known {
		class = ` class="language-` + html.EscapeString(fence.language) + `"`
	}
	var out strings.Builder
	out.WriteString(`<div class="code-block">`)
	if fence.language != "" || fence.filename != "" {
		out.WriteString(`<div class="code-header">`)
		if fence.language != "" {
			out.WriteString(`<span class="code-language">`)
			out.WriteString(html.EscapeString(a1LanguageLabel(fence.language)))
			out.WriteString(`</span>`)
		}
		if fence.filename != "" {
			out.WriteString(`<span class="code-filename">`)
			out.WriteString(html.EscapeString(fence.filename))
			out.WriteString(`</span>`)
		}
		// The button is hidden until a future first-party enhancement is
		// available. With JavaScript disabled it is absent as a control.
		out.WriteString(`<button type="button" class="code-copy" data-copy-code="`)
		out.WriteString(html.EscapeString(source))
		out.WriteString(`" aria-label="Copy code" hidden>Copy</button>`)
		out.WriteString(`</div>`)
	}
	out.WriteString(`<pre><code`)
	out.WriteString(class)
	out.WriteByte('>')
	out.WriteString(code)
	out.WriteString(`</code></pre></div>` + "\n")
	return out.String()
}

func a1LanguageLabel(language string) string {
	switch strings.ToLower(language) {
	case "go", "golang":
		return "Go"
	case "bash", "sh", "shell":
		return "Bash"
	case "templ":
		return "templ"
	case "json":
		return "JSON"
	case "yaml", "yml":
		return "YAML"
	default:
		return language
	}
}

func seriesNavigationFor(post Post, posts []Post) a1SeriesNav {
	if post.Series == "" {
		return a1SeriesNav{}
	}
	members := make([]Post, 0)
	for _, candidate := range posts {
		if candidate.Series == post.Series {
			members = append(members, candidate)
		}
	}
	if len(members) < 2 {
		return a1SeriesNav{}
	}
	sort.SliceStable(members, func(i, j int) bool {
		leftPart, leftErr := strconv.Atoi(strings.TrimSpace(members[i].Part))
		rightPart, rightErr := strconv.Atoi(strings.TrimSpace(members[j].Part))
		if leftErr == nil && rightErr == nil && leftPart != rightPart {
			return leftPart < rightPart
		}
		if members[i].Part != members[j].Part {
			return members[i].Part < members[j].Part
		}
		if members[i].Date != members[j].Date {
			return members[i].Date < members[j].Date
		}
		return members[i].Slug < members[j].Slug
	})
	index := -1
	for i := range members {
		if members[i].Slug == post.Slug {
			index = i
			break
		}
	}
	if index < 0 {
		return a1SeriesNav{}
	}
	navigation := a1SeriesNav{
		Series:   post.Series,
		Position: fmt.Sprintf("Part %d of %d", index+1, len(members)),
	}
	if index > 0 {
		previous := members[index-1]
		navigation.Previous = &previous
		navigation.PreviousText = previous.Title
	}
	if index+1 < len(members) {
		next := members[index+1]
		navigation.Next = &next
		navigation.NextText = next.Title
	}
	return navigation
}
