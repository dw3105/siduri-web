package site

import (
	"html"
	"strings"
	"unicode"
)

// a1KnownLanguage is deliberately closed. Unknown fences are still emitted
// as escaped code, so a post cannot fail merely because its language label is
// new to this small build-time tokenizer.
func a1KnownLanguage(language string) bool {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "go", "golang", "bash", "sh", "shell", "templ", "json", "yaml", "yml":
		return true
	default:
		return false
	}
}

func a1HighlightCode(language, source string) string {
	language = strings.ToLower(strings.TrimSpace(language))
	if !a1KnownLanguage(language) {
		return html.EscapeString(source)
	}
	return a1HighlightTokens(language, source)
}

type a1CodeToken struct {
	class string
	text  string
}

func a1HighlightTokens(language, source string) string {
	tokens := make([]a1CodeToken, 0, len(source)/4)
	for index := 0; index < len(source); {
		if isA1Space(source[index]) {
			start := index
			for index < len(source) && isA1Space(source[index]) {
				index++
			}
			tokens = append(tokens, a1CodeToken{text: source[start:index]})
			continue
		}

		if a1StartsComment(language, source[index:]) {
			start := index
			for index < len(source) && source[index] != '\n' {
				index++
			}
			tokens = append(tokens, a1CodeToken{class: "comment", text: source[start:index]})
			continue
		}

		if strings.HasPrefix(source[index:], "/*") {
			start := index
			index += 2
			if end := strings.Index(source[index:], "*/"); end >= 0 {
				index += end + 2
			} else {
				index = len(source)
			}
			tokens = append(tokens, a1CodeToken{class: "comment", text: source[start:index]})
			continue
		}

		if quote := a1QuoteAt(language, source[index]); quote != 0 {
			start := index
			index = a1QuotedEnd(source, index, quote)
			class := "string"
			if (language == "json" || language == "yaml" || language == "yml") && a1NextNonSpace(source, index) == ':' {
				class = "property"
			}
			tokens = append(tokens, a1CodeToken{class: class, text: source[start:index]})
			continue
		}

		if unicode.IsDigit(rune(source[index])) {
			start := index
			for index < len(source) && (unicode.IsDigit(rune(source[index])) || strings.ContainsRune("._", rune(source[index]))) {
				index++
			}
			tokens = append(tokens, a1CodeToken{class: "number", text: source[start:index]})
			continue
		}

		if a1IdentifierStart(source[index]) {
			start := index
			for index < len(source) && a1IdentifierPart(language, source[index]) {
				index++
			}
			word := source[start:index]
			class := a1WordClass(language, word, source, index)
			tokens = append(tokens, a1CodeToken{class: class, text: word})
			continue
		}

		if source[index] == '$' && index+1 < len(source) && a1IdentifierStart(source[index+1]) {
			start := index
			index += 2
			for index < len(source) && a1IdentifierPart(language, source[index]) {
				index++
			}
			tokens = append(tokens, a1CodeToken{class: "variable", text: source[start:index]})
			continue
		}

		tokens = append(tokens, a1CodeToken{text: source[index : index+1]})
		index++
	}

	var out strings.Builder
	for _, token := range tokens {
		text := html.EscapeString(token.text)
		if token.class == "" {
			out.WriteString(text)
		} else {
			out.WriteString(`<span class="tok-`)
			out.WriteString(token.class)
			out.WriteString(`">`)
			out.WriteString(text)
			out.WriteString(`</span>`)
		}
	}
	return out.String()
}

func isA1Space(value byte) bool {
	return value == ' ' || value == '\t' || value == '\n' || value == '\r'
}

func a1StartsComment(language, source string) bool {
	if (language == "go" || language == "golang" || language == "templ") && strings.HasPrefix(source, "//") {
		return true
	}
	return (language == "bash" || language == "sh" || language == "shell" || language == "yaml" || language == "yml") && source[0] == '#'
}

func a1QuoteAt(language string, value byte) byte {
	switch value {
	case '"', '\'', '`':
		if language == "json" && value != '"' {
			return 0
		}
		return value
	default:
		return 0
	}
}

func a1QuotedEnd(source string, start int, quote byte) int {
	for index := start + 1; index < len(source); index++ {
		if source[index] == '\\' && quote != '`' {
			index++
			continue
		}
		if source[index] == quote {
			return index + 1
		}
	}
	return len(source)
}

func a1NextNonSpace(source string, index int) byte {
	for index < len(source) && isA1Space(source[index]) {
		index++
	}
	if index < len(source) {
		return source[index]
	}
	return 0
}

func a1IdentifierStart(value byte) bool {
	return value == '_' || unicode.IsLetter(rune(value))
}

func a1IdentifierPart(language string, value byte) bool {
	if a1IdentifierStart(value) || unicode.IsDigit(rune(value)) {
		return true
	}
	return (language == "yaml" || language == "yml") && value == '-'
}

func a1WordClass(language, word, source string, end int) string {
	if a1A1Keywords(language)[word] {
		return "keyword"
	}
	if word == "true" || word == "false" || word == "null" || word == "nil" {
		return "boolean"
	}
	if (language == "json" || language == "yaml" || language == "yml") && a1NextNonSpace(source, end) == ':' {
		return "property"
	}
	if a1NextNonSpace(source, end) == '(' {
		return "function"
	}
	return ""
}

func a1A1Keywords(language string) map[string]bool {
	keywords := map[string]bool{}
	switch language {
	case "go", "golang", "templ":
		for _, word := range strings.Fields("break default func interface select case defer go map struct chan else goto package switch const fallthrough if range type continue for import return var") {
			keywords[word] = true
		}
	case "bash", "sh", "shell":
		for _, word := range strings.Fields("if then elif else fi for while in do done case esac function select") {
			keywords[word] = true
		}
	case "json", "yaml", "yml":
		// JSON and YAML have no keyword set beyond their scalar values.
	}
	return keywords
}
