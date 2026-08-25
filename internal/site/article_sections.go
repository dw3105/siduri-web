package site

import (
	"fmt"
	"sort"
	"sync"

	"github.com/a-h/templ"
)

// ArticleSectionSlot is the closed set of positions a lane may extend in an
// article page. The components themselves own any markup they need.
type ArticleSectionSlot string

const (
	ArticleSectionBodyAside ArticleSectionSlot = "body-aside"
	ArticleSectionAfterBody ArticleSectionSlot = "after-body"
	ArticleSectionFooter    ArticleSectionSlot = "footer"
)

// ArticleSection is one additive contribution to an article page. Names are
// sorted within each slot before rendering, so registration order is not part
// of the output.
type ArticleSection struct {
	Name   string
	Slot   ArticleSectionSlot
	Render func(Post, PageData) templ.Component
}

// ArticleSections is the set of rendered article contributions passed to the
// article template. It is build-local even though registrations are declared
// at package initialization time.
type ArticleSections struct {
	BodyAside []templ.Component
	AfterBody []templ.Component
	Footer    []templ.Component
}

var (
	articleSectionMu sync.RWMutex
	articleSections  = make(map[string]ArticleSection)
	bodyRendererMu   sync.RWMutex
	bodyRenderer     *BodyRenderer
)

// RegisterArticleSection lets a lane add one component to a named article
// slot. Duplicate names and unknown slots are programmer errors and panic.
func RegisterArticleSection(section ArticleSection) {
	if section.Name == "" || section.Render == nil {
		panic("site: incomplete article section registration")
	}
	if !validArticleSectionSlot(section.Slot) {
		panic(fmt.Sprintf("site: unknown article section slot %q", section.Slot))
	}
	articleSectionMu.Lock()
	defer articleSectionMu.Unlock()
	if _, exists := articleSections[section.Name]; exists {
		panic(fmt.Sprintf("site: duplicate article section %q", section.Name))
	}
	articleSections[section.Name] = section
}

func validArticleSectionSlot(slot ArticleSectionSlot) bool {
	switch slot {
	case ArticleSectionBodyAside, ArticleSectionAfterBody, ArticleSectionFooter:
		return true
	default:
		return false
	}
}

func articleSectionsFor(post Post, data PageData) ArticleSections {
	articleSectionMu.RLock()
	registered := make([]ArticleSection, 0, len(articleSections))
	for _, section := range articleSections {
		registered = append(registered, section)
	}
	articleSectionMu.RUnlock()
	sort.Slice(registered, func(i, j int) bool {
		if registered[i].Slot != registered[j].Slot {
			return registered[i].Slot < registered[j].Slot
		}
		return registered[i].Name < registered[j].Name
	})

	sections := ArticleSections{}
	for _, section := range registered {
		component := section.Render(post, data)
		if component == nil {
			panic(fmt.Sprintf("site: article section %q returned nil", section.Name))
		}
		switch section.Slot {
		case ArticleSectionBodyAside:
			sections.BodyAside = append(sections.BodyAside, component)
		case ArticleSectionAfterBody:
			sections.AfterBody = append(sections.AfterBody, component)
		case ArticleSectionFooter:
			sections.Footer = append(sections.Footer, component)
		}
	}
	return sections
}

// BodyRenderer owns the Markdown-to-HTML step for article pages. There may be
// only one owner; a missing owner uses renderMarkdown.
type BodyRenderer struct {
	Name   string
	Render func(Post, PageData) string
}

// RegisterBodyRenderer replaces the default article body renderer. A second
// registration is a panic because body rendering is ownership, not a slot.
func RegisterBodyRenderer(renderer BodyRenderer) {
	if renderer.Name == "" || renderer.Render == nil {
		panic("site: incomplete body renderer registration")
	}
	bodyRendererMu.Lock()
	defer bodyRendererMu.Unlock()
	if bodyRenderer != nil {
		panic(fmt.Sprintf("site: duplicate body renderer %q (already owned by %q)", renderer.Name, bodyRenderer.Name))
	}
	copy := renderer
	bodyRenderer = &copy
}

func renderArticleBody(post Post, data PageData) string {
	bodyRendererMu.RLock()
	renderer := bodyRenderer
	bodyRendererMu.RUnlock()
	if renderer == nil {
		return renderMarkdown(post.Body)
	}
	return renderer.Render(post, data)
}
