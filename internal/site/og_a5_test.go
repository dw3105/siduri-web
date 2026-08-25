package site

import (
	"bytes"
	"image/png"
	"testing"
)

func TestA5OGImageIsDeterministicAndSized(t *testing.T) {
	post := Post{Title: "A deterministic card", Series: "Build log", Date: "2026-08-25"}
	first := renderOGImage(post)
	second := renderOGImage(post)
	if !bytes.Equal(first, second) {
		t.Fatal("identical post metadata produced different OG image bytes")
	}
	decoded, err := png.Decode(bytes.NewReader(first))
	if err != nil {
		t.Fatalf("decode OG image: %v", err)
	}
	if got := decoded.Bounds().Size(); got.X != ogWidth || got.Y != ogHeight {
		t.Fatalf("OG image size = %v, want %dx%d", got, ogWidth, ogHeight)
	}
	variants := []Post{
		{Title: "A deterministic card", Series: "A different series", Date: "2026-08-25"},
		{Title: "A deterministic card", Series: "Build log", Date: "2026-08-26"},
		{Title: "A different card", Series: "Build log", Date: "2026-08-25"},
	}
	for _, variant := range variants {
		if bytes.Equal(first, renderOGImage(variant)) {
			t.Fatalf("changing post metadata did not change OG image bytes: %#v", variant)
		}
	}
}
