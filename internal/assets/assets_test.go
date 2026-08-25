package assets

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessImageHashUsesProcessedBytes(t *testing.T) {
	source := testPNG(t, 12, 7, color.RGBA{R: 20, G: 40, B: 60, A: 255})
	first, err := ProcessImage("assets/example.png", source)
	if err != nil {
		t.Fatal(err)
	}
	second, err := ProcessImage("assets/example.png", source)
	if err != nil {
		t.Fatal(err)
	}
	if first.OutputPath != second.OutputPath || !bytes.Equal(first.Bytes, second.Bytes) {
		t.Fatalf("same source was not deterministic: %q and %q", first.OutputPath, second.OutputPath)
	}

	changedSource := testPNG(t, 12, 7, color.RGBA{R: 21, G: 40, B: 60, A: 255})
	changed, err := ProcessImage("assets/example.png", changedSource)
	if err != nil {
		t.Fatal(err)
	}
	if changed.OutputPath == first.OutputPath {
		t.Fatalf("changed source kept content hash %q", changed.OutputPath)
	}
	if !strings.Contains(first.OutputPath, first.Hash[:contentHashLength]) {
		t.Fatalf("output path does not contain processed-byte hash: %q", first.OutputPath)
	}
}

func TestDiscoverRejectsUnsupportedImageByName(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "assets", "camera.webp"), []byte("not an encoder input"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Discover(root)
	if err == nil || !strings.Contains(err.Error(), "camera.webp") {
		t.Fatalf("expected named unsupported image error, got %v", err)
	}
}

func TestImageTagCarriesExplicitDimensions(t *testing.T) {
	image := Image{URL: "/assets/example-abc.png", Width: 12, Height: 7}
	tag := ImageTag(image, `A "picture"`, true)
	for _, required := range []string{
		`src="/assets/example-abc.png"`,
		`alt="A &#34;picture&#34;"`,
		`width="12"`,
		`height="7"`,
		`loading="lazy"`,
	} {
		if !strings.Contains(tag, required) {
			t.Fatalf("image tag missing %s: %s", required, tag)
		}
	}
}

func testPNG(t *testing.T, width, height int, fill color.Color) []byte {
	t.Helper()
	image := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			image.Set(x, y, fill)
		}
	}
	var data bytes.Buffer
	if err := png.Encode(&data, image); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}
