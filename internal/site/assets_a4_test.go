package site

import (
	"bytes"
	"compress/gzip"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/dw3105/siduri-web/internal/assets"
)

func TestA4ImageTagsHaveDimensionsAndLazyBelowFold(t *testing.T) {
	first := testA4Image(t, "assets/first.png", 640, 360)
	second := testA4Image(t, "assets/second.png", 320, 180)
	images := map[string]assets.Image{first.Source: first, second.Source: second}

	result, err := assets.RewriteMarkdownImages("![First](assets/first.png)\n\n![Second](/assets/second.png)", images)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(result, "<img ") != 2 {
		t.Fatalf("expected two image tags, got %q", result)
	}
	if !strings.Contains(result, `width="640" height="360"`) || !strings.Contains(result, `width="320" height="180"`) {
		t.Fatalf("image dimensions missing: %q", result)
	}
	if strings.Count(result, `loading="lazy"`) != 1 {
		t.Fatalf("expected only below-fold image to be lazy: %q", result)
	}
	if strings.Contains(result, "assets/first.png") || strings.Contains(result, "assets/second.png") {
		t.Fatalf("source image URL leaked into output: %q", result)
	}
}

func TestA4UnprocessedImageNamesFile(t *testing.T) {
	_, err := assets.RewriteMarkdownImages("![Missing](/assets/missing.png)", nil)
	if err == nil || !strings.Contains(err.Error(), `missing.png`) {
		t.Fatalf("expected named unprocessed image error, got %v", err)
	}
}

func TestA4CSSBudget(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "static", "site.css"))
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, criticalCSS...)
	var compressed bytes.Buffer
	writer, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	t.Logf("CSS gzip: %d bytes (budget: 15360 bytes)", compressed.Len())
	if compressed.Len() > 15*1024 {
		t.Fatalf("CSS gzip budget exceeded: %d bytes", compressed.Len())
	}
}

func TestA4HeadHasAtMostTwoFontPreloads(t *testing.T) {
	component := Document("A4", false, renderHeadFragments(PageData{}), templ.Raw("content"))
	data := renderComponent(t, component)
	preloadPattern := regexp.MustCompile(`<link rel="preload"[^>]+as="font"`)
	if count := len(preloadPattern.FindAll(data, -1)); count > 2 {
		t.Fatalf("expected at most two font preloads, got %d", count)
	}
}

func TestA4AssetRouteUsesByteOutputs(t *testing.T) {
	set := newRouteSet()
	registerA4Assets(PageData{}, set)
	for _, route := range set.routes {
		for _, output := range route.Output.Expand(PageData{}) {
			if output.Render != nil || output.Bytes == nil {
				t.Fatalf("asset output is not a ByteOutput: %#v", output)
			}
		}
	}
}

func testA4Image(t *testing.T, source string, width, height int) assets.Image {
	t.Helper()
	image, err := assets.ProcessImage(source, testA4PNGBytes(t, width, height))
	if err != nil {
		t.Fatal(err)
	}
	return image
}

func testA4PNGBytes(t *testing.T, width, height int) []byte {
	t.Helper()
	var data bytes.Buffer
	image := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			image.SetRGBA(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: 120, A: 255})
		}
	}
	if err := png.Encode(&data, image); err != nil {
		t.Fatal(err)
	}
	return data.Bytes()
}
