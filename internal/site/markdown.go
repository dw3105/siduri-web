package site

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var linkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^\s)]+)\)`)

var (
	strongPattern      = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	underStrongPattern = regexp.MustCompile(`__([^_]+)__`)
	codePattern        = regexp.MustCompile("`([^`]+)`")
	emphasisPattern    = regexp.MustCompile(`\*([^*]+)\*`)
)

func renderMarkdown(source string) string {
	lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")
	var out strings.Builder
	paragraph := make([]string, 0)
	inCode := false
	var codeLanguage string
	flushParagraph := func() {
		if len(paragraph) == 0 {
			return
		}
		out.WriteString("<p>")
		out.WriteString(renderInline(strings.Join(paragraph, "\n")))
		out.WriteString("</p>\n")
		paragraph = paragraph[:0]
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if inCode {
				out.WriteString("</code></pre>\n")
				inCode = false
				codeLanguage = ""
			} else {
				flushParagraph()
				codeLanguage = strings.TrimSpace(strings.TrimPrefix(line, "```"))
				class := ""
				if codeLanguage != "" {
					class = ` class="language-` + html.EscapeString(codeLanguage) + `"`
				}
				out.WriteString("<pre><code" + class + ">")
				inCode = true
			}
			continue
		}
		if inCode {
			out.WriteString(html.EscapeString(line))
			out.WriteByte('\n')
			continue
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			flushParagraph()
			continue
		}
		if strings.HasPrefix(trimmed, "### ") {
			flushParagraph()
			out.WriteString("<h3>")
			out.WriteString(renderInline(strings.TrimPrefix(trimmed, "### ")))
			out.WriteString("</h3>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			flushParagraph()
			out.WriteString("<h2>")
			out.WriteString(renderInline(strings.TrimPrefix(trimmed, "## ")))
			out.WriteString("</h2>\n")
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			flushParagraph()
			out.WriteString("<h2>")
			out.WriteString(renderInline(strings.TrimPrefix(trimmed, "# ")))
			out.WriteString("</h2>\n")
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
		out.WriteString("</code></pre>\n")
	} else {
		flushParagraph()
	}
	return out.String()
}

func renderInline(source string) string {
	const marker = "\x00LINK%d\x00"
	links := make([]string, 0)
	protected := linkPattern.ReplaceAllStringFunc(source, func(match string) string {
		parts := linkPattern.FindStringSubmatch(match)
		url := parts[2]
		if !safeURL(url) {
			return html.EscapeString(parts[1])
		}
		links = append(links, `<a href="`+html.EscapeString(url)+`">`+renderInline(parts[1])+`</a>`)
		return formatMarker(marker, len(links)-1)
	})
	escaped := html.EscapeString(protected)
	escaped = strongPattern.ReplaceAllString(escaped, `<strong>$1</strong>`)
	escaped = underStrongPattern.ReplaceAllString(escaped, `<strong>$1</strong>`)
	escaped = codePattern.ReplaceAllString(escaped, `<code>$1</code>`)
	escaped = emphasisPattern.ReplaceAllString(escaped, `<em>$1</em>`)
	for i, link := range links {
		escaped = strings.ReplaceAll(escaped, html.EscapeString(formatMarker(marker, i)), link)
	}
	return strings.ReplaceAll(escaped, "\n", "<br />\n")
}

func formatMarker(format string, n int) string {
	return fmt.Sprintf(format, n)
}

func safeURL(url string) bool {
	return strings.HasPrefix(url, "/") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")
}
