package site

import (
	"strings"
	"testing"
)

func TestA1HighlighterCoversSiteLanguages(t *testing.T) {
	cases := map[string]string{
		"go":    "package main\nfunc main() {}",
		"bash":  "if test -n \"$HOME\"; then\n  echo yes\nfi",
		"templ": "templ page() { <p>{ value }</p> }",
		"json":  "{\"name\": true, \"count\": 2}",
		"yaml":  "name: siduri\ncount: 2",
	}
	for language, source := range cases {
		language, source := language, source
		t.Run(language, func(t *testing.T) {
			rendered := a1HighlightCode(language, source)
			if !strings.Contains(rendered, `<span class="tok-`) {
				t.Fatalf("%s produced no token spans: %s", language, rendered)
			}
			if strings.Contains(rendered, source) && strings.ContainsAny(source, "<>&\"") {
				t.Fatalf("%s was not safely escaped: %s", language, rendered)
			}
		})
	}
}

func TestA1HighlighterUnknownLanguageIsEscapedPlainText(t *testing.T) {
	rendered := a1HighlightCode("rust", "fn main() { <panic> }")
	if strings.Contains(rendered, `<span`) {
		t.Fatal("unknown language unexpectedly produced token spans")
	}
	if rendered != "fn main() { &lt;panic&gt; }" {
		t.Fatalf("unknown language was not escaped plainly: %q", rendered)
	}
}
