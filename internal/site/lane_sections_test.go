package site

import (
	"bytes"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

func TestArticleSectionsSortWithinSlotsAndIgnoreRegistrationOrder(t *testing.T) {
	post := Post{
		Title:        "Section fixture",
		Date:         "2026-08-25",
		PlainSummary: "A section fixture.",
		Tags:         []string{"method"},
		Body:         "The body.",
	}
	sections := []ArticleSection{
		{
			Name: "lane-zeta",
			Slot: ArticleSectionBodyAside,
			Render: func(Post, PageData) templ.Component {
				return templ.Raw(`<aside data-section="zeta">zeta</aside>`)
			},
		},
		{
			Name: "lane-alpha",
			Slot: ArticleSectionBodyAside,
			Render: func(Post, PageData) templ.Component {
				return templ.Raw(`<aside data-section="alpha">alpha</aside>`)
			},
		},
		{
			Name: "lane-after",
			Slot: ArticleSectionAfterBody,
			Render: func(Post, PageData) templ.Component {
				return templ.Raw(`<section data-section="after">after</section>`)
			},
		},
	}

	first := renderArticleWithSections(t, post, sections)
	second := renderArticleWithSections(t, post, []ArticleSection{sections[2], sections[1], sections[0]})
	if !bytes.Equal(first, second) {
		t.Fatalf("registration order changed article output\nfirst:  %s\nsecond: %s", first, second)
	}

	output := string(first)
	for _, marker := range []string{`data-section="alpha"`, `data-section="zeta"`, `data-section="after"`} {
		if !strings.Contains(output, marker) {
			t.Fatalf("article section %q missing from output", marker)
		}
	}
	if strings.Index(output, `data-section="alpha"`) > strings.Index(output, `data-section="zeta"`) {
		t.Fatal("same-slot sections were not sorted by name")
	}
	if strings.Index(output, `data-section="after"`) < strings.Index(output, `</div>`) {
		t.Fatal("after-body section rendered before the body")
	}
}

func renderArticleWithSections(t *testing.T, post Post, registrations []ArticleSection) []byte {
	t.Helper()
	restore := isolateArticleSections()
	defer restore()
	for _, registration := range registrations {
		RegisterArticleSection(registration)
	}
	return renderComponent(t, ArticlePage(post, renderMarkdown(post.Body), false))
}

func isolateArticleSections() func() {
	articleSectionMu.Lock()
	saved := articleSections
	articleSections = make(map[string]ArticleSection)
	articleSectionMu.Unlock()
	return func() {
		articleSectionMu.Lock()
		articleSections = saved
		articleSectionMu.Unlock()
	}
}
