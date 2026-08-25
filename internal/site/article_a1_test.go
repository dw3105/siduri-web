package site

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestA1TableOfContentsUsesFixtureThreshold(t *testing.T) {
	short := Post{
		Title:        "Short fixture",
		Slug:         "short-fixture",
		Date:         "2026-08-25",
		PlainSummary: "Short fixture.",
		Body:         "## A heading\n\n" + strings.Repeat("word ", a1TOCWordThreshold-10),
	}
	long := short
	long.Title = "Long fixture"
	long.Slug = "long-fixture"
	long.Body = "## A heading\n\n" + strings.Repeat("word ", a1TOCWordThreshold+1)

	shortHTML := string(renderComponent(t, ArticlePage(short, renderA1Markdown(short.Body), false)))
	longHTML := string(renderComponent(t, ArticlePage(long, renderA1Markdown(long.Body), false)))
	if strings.Contains(shortHTML, `class="article-toc"`) {
		t.Fatal("short fixture unexpectedly rendered a table of contents")
	}
	if !strings.Contains(longHTML, `class="article-toc"`) {
		t.Fatal("long fixture did not render a table of contents")
	}
	if !strings.Contains(longHTML, `href="#a-heading"`) || !strings.Contains(longHTML, `id="a-heading"`) {
		t.Fatal("long fixture table of contents is not linked to the generated heading")
	}
}

func TestA1HighlightingAndUnknownLanguageFallback(t *testing.T) {
	source := "```go filename=main.go\npackage main\n\nfunc main() {\n\tprintln(\"hi\")\n}\n```\n\n```nosuchlang\n<tag>\n```"
	rendered := renderA1Markdown(source)
	if !strings.Contains(rendered, `<span class="tok-keyword">package</span>`) {
		t.Fatal("Go keyword was not highlighted at build time")
	}
	if !strings.Contains(rendered, `<span class="tok-string">&#34;hi&#34;</span>`) {
		t.Fatal("Go string was not highlighted at build time")
	}
	if !strings.Contains(rendered, `class="code-language">Go</span>`) {
		t.Fatal("code language label is missing")
	}
	if !strings.Contains(rendered, `class="code-filename">main.go</span>`) {
		t.Fatal("code filename caption is missing")
	}
	if !strings.Contains(rendered, `aria-label="Copy code" hidden`) {
		t.Fatal("copy control is not disabled by default")
	}
	unknownStart := strings.Index(rendered, `<pre><code>`)
	if unknownStart < 0 {
		t.Fatal("unknown language did not fall back to a plain pre/code block")
	}
	unknownEnd := strings.Index(rendered[unknownStart:], `</code>`)
	if unknownEnd < 0 || strings.Contains(rendered[unknownStart:unknownStart+unknownEnd], `<span`) {
		t.Fatal("unknown language was transformed instead of rendered plainly")
	}
	if !strings.Contains(rendered, `&lt;tag&gt;`) {
		t.Fatal("unknown language code was not HTML escaped")
	}
}

func TestA1SeriesNavigationUsesPartOrder(t *testing.T) {
	posts := []Post{
		{Title: "Part three", Slug: "part-three", Date: "2026-08-27", Series: "Build diary", Part: "3"},
		{Title: "Part one", Slug: "part-one", Date: "2026-08-25", Series: "Build diary", Part: "1"},
		{Title: "Part two", Slug: "part-two", Date: "2026-08-26", Series: "Build diary", Part: "2"},
	}
	navigation := seriesNavigationFor(posts[0], posts)
	if navigation.Position != "Part 3 of 3" {
		t.Fatalf("unexpected series position: %q", navigation.Position)
	}
	if navigation.Previous == nil || navigation.Previous.Slug != "part-two" {
		t.Fatalf("previous series member is wrong: %#v", navigation.Previous)
	}
	if navigation.Next != nil {
		t.Fatalf("last series member unexpectedly has a next link: %#v", navigation.Next)
	}

	rendered := string(renderComponent(t, a1SeriesNavigation(seriesNavigationFor(posts[1], posts))))
	if !strings.Contains(rendered, "Part 1 of 3") || !strings.Contains(rendered, "/journal/part-two/") {
		t.Fatalf("first series member did not render the next link: %s", rendered)
	}
	withoutSeries := string(renderComponent(t, a1SeriesNavigation(seriesNavigationFor(Post{Slug: "standalone"}, posts))))
	if strings.TrimSpace(withoutSeries) != "" {
		t.Fatalf("post without a series rendered a navigator: %q", withoutSeries)
	}
}

func TestA1BodyRendererIsBuildTimeAndKeepsPlainMarkdown(t *testing.T) {
	post := Post{Body: "A paragraph.\n\n```bash\necho $HOME\n```"}
	rendered := renderArticleBody(post, PageData{})
	if !strings.Contains(rendered, `<span class="tok-variable">$HOME</span>`) {
		t.Fatal("registered A1 body renderer did not tokenize a fenced block")
	}
	if !strings.Contains(rendered, `<p>A paragraph.</p>`) {
		t.Fatal("A1 body renderer dropped ordinary Markdown paragraphs")
	}
	if strings.Contains(rendered, "<script") {
		t.Fatal("article body introduced a client-side script")
	}
}

func TestA1StylesheetIsAByteOutput(t *testing.T) {
	set := newRouteSet()
	registerA1Assets(PageData{}, set)
	route, ok := set.routes["a1-article-stylesheet"]
	if !ok {
		t.Fatal("A1 stylesheet route was not registered")
	}
	outputs := route.Output.Expand(PageData{})
	if len(outputs) != 1 || outputs[0].Bytes == nil || outputs[0].Render != nil {
		t.Fatalf("stylesheet was not emitted as one byte output: %#v", outputs)
	}
	want, err := os.ReadFile(filepath.Join("..", "..", "static", "article_a1.css"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(outputs[0].Bytes, want) {
		t.Fatal("stylesheet byte output differs from the owned static asset")
	}
}
